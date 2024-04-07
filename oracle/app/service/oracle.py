import datetime
import uuid

from app.api.oracle import oracle_pb2_grpc, oracle_pb2
from app.db.pg_repo import OracleRepository
from app.houston import loggy

import pandas as pd
from lightfm.data import Dataset
from lightfm import LightFM


class OracleService(oracle_pb2_grpc.OracleServicer):
    def __init__(self, repo: OracleRepository) -> None:
        self.__repo = repo

    def RetrainLFMv1(self, request, context):
        LEARNING_RATE = 0.25
        NO_EPOCHS = 20
        NO_COMPONENTS = 20
        ITEM_ALPHA = 1e-6
        USER_ALPHA = 1e-6

        loggy.info("start retrain LFMv1")
        start = datetime.datetime.now()

        data = pd.DataFrame(self.__repo.get_lfmv1_train_data(),
                            columns=['from', 'to', 'positive', 'updated_at__from', 'updated_at__to'])
        data['rating'] = data['positive'] * 1

        from_all_feature_list = []
        from_all_feature_list.extend(list(data["updated_at__from"].drop_duplicates()))

        to_all_feature_list = []
        to_all_feature_list.extend(list(data["updated_at__to"].drop_duplicates()))

        dataset = Dataset()
        dataset.fit(users=data["from"], items=data["to"], user_features=from_all_feature_list,
                    item_features=to_all_feature_list)

        data_with_unique_from = data.drop_duplicates("from")
        data_with_unique_to = data.drop_duplicates("to")

        # TODO: add features
        from_features_list = [(user_id, [updated_at]) for user_id, updated_at in
                              zip(data_with_unique_from["from"], data_with_unique_from["updated_at__from"])]

        to_features_list = [(user_id, [updated_at]) for user_id, updated_at in
                            zip(data_with_unique_to["to"], data_with_unique_to["updated_at__to"])]

        from_features = dataset.build_user_features(from_features_list)
        to_features = dataset.build_item_features(to_features_list)

        interactions, weights = dataset.build_interactions(data[["from", "to", "rating"]].values)

        user_id_map, user_feature_map, item_id_map, feature_item_map = dataset.mapping()

        model = LightFM(loss="warp",
                        no_components=NO_COMPONENTS,
                        learning_rate=LEARNING_RATE,
                        item_alpha=ITEM_ALPHA,
                        user_alpha=USER_ALPHA)

        model.fit(interactions=interactions,
                  user_features=from_features,
                  item_features=to_features,
                  epochs=NO_EPOCHS)

        all_user_ids = data_with_unique_from["from"]
        all_scores = []

        for user_id in all_user_ids:
            list_scores = model.predict(user_id_map[user_id], list(item_id_map.values()))

            for pair, score in zip(item_id_map.keys(), list_scores):
                all_scores.append((uuid.uuid4(), user_id, pair, score))

        self.__repo.update_lfmv1_results(all_scores)

        finish = datetime.datetime.now()
        loggy.info(f'retrain LFMv1 elapsed: {finish - start}')

        return oracle_pb2.RetrainLFMv1Response()

    def GetTaskStatus(self, request, context):
        loggy.info("stub for future release")

        return oracle_pb2.GetTaskStatusResponse()
