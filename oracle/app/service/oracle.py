import datetime
import uuid

import grpc

from app.api.oracle import oracle_pb2_grpc, oracle_pb2
from app.db.pg_repo import OracleRepository
from app.houston import loggy

import pandas as pd
from lightfm.data import Dataset
from lightfm import LightFM


class OracleService(oracle_pb2_grpc.OracleServicer):
    def __init__(self, repo: OracleRepository, lfmv1_config) -> None:
        self.__repo = repo
        self.__lfmv1_config = lfmv1_config

    def RetrainLFMv1(self, request, context):
        loggy.info("start retrain LFMv1")
        start = datetime.datetime.now()

        from_feature_names = [
            'interests__from', 'birthday__from', 'gender__from', 'location_lat__from', 'location_lon__from'
        ]

        to_feature_names = [
            'interests__to', 'birthday__to', 'gender__to', 'location_lat__to', 'location_lon__to'
        ]

        cols = ['from', 'to', 'positive']
        for from_feature, to_feature in zip(from_feature_names, to_feature_names):
            cols.append(from_feature)
            cols.append(to_feature)

        try:
            data = pd.DataFrame(self.__repo.get_lfmv1_train_data(), columns=cols)
            data['rating'] = data['positive'] * 2 - 1

            data['interests__from'] = data['interests__from'].apply(
                lambda y: [str(x) + '_interest__from' for x in y or []])

            data['interests__to'] = data['interests__to'].apply(lambda y: [str(x) + '_interest__to' for x in y or []])

            data = data.drop('interests__from', axis=1).join(data['interests__from'].str.join('|').str.get_dummies())
            data = data.drop('interests__to', axis=1).join(data['interests__to'].str.join('|').str.get_dummies())

            from_feature_names = [f_name for f_name in data.columns if f_name.endswith('__from')]
            to_feature_names = [f_name for f_name in data.columns if f_name.endswith('__to')]

            from_all_feature_list = []
            for feat in from_feature_names:
                from_all_feature_list.extend(list(data[feat]))

            to_all_feature_list = []
            for feat in to_feature_names:
                to_all_feature_list.extend(list(data[feat]))


            dataset = Dataset()
            dataset.fit(users=data["from"], items=data["to"], user_features=from_all_feature_list,
                        item_features=to_all_feature_list)

            data_with_unique_from = data.drop_duplicates("from")
            data_with_unique_to = data.drop_duplicates("to")

            from_features_list = []

            for id, (_, row) in zip(data_with_unique_from['from'],
                                    data_with_unique_from[from_feature_names].iterrows()):
                l = []
                for feat in row:
                    l.append(feat)

                from_features_list.append((id, l))

            to_features_list = []

            for id, (_, row) in zip(data_with_unique_to['to'], data_with_unique_to[to_feature_names].iterrows()):
                l = []
                for feat in row:
                    l.append(feat)

                to_features_list.append((id, l))

            from_features = dataset.build_user_features(from_features_list)
            to_features = dataset.build_item_features(to_features_list)

            interactions, weights = dataset.build_interactions(data[["from", "to", "rating"]].values)
            user_id_map, user_feature_map, item_id_map, feature_item_map = dataset.mapping()

            model = LightFM(loss="warp",
                            no_components=self.__lfmv1_config.get("no_components"),
                            learning_rate=self.__lfmv1_config.get("learning_rate"),
                            item_alpha=self.__lfmv1_config.get("item_alpha"),
                            user_alpha=self.__lfmv1_config.get("user_alpha"))

            model.fit(interactions=interactions,
                      user_features=from_features,
                      item_features=to_features,
                      epochs=self.__lfmv1_config.get("no_epochs"))

            all_user_ids = data_with_unique_from["from"]
            all_scores = []

            for user_id in all_user_ids:
                list_scores = model.predict(user_id_map[user_id], list(item_id_map.values()))

                for pair, score in zip(item_id_map.keys(), list_scores):
                    all_scores.append((uuid.uuid4(), user_id, pair, score))

            self.__repo.update_lfmv1_results(all_scores)
        except Exception as ex:
            loggy.info(f"Failed to retrain LFMv1 model: {ex}")
            return oracle_pb2.RetrainLFMv1Response()

        finish = datetime.datetime.now()
        loggy.info(f'retrain LFMv1 elapsed: {finish - start}')

        return oracle_pb2.RetrainLFMv1Response()

    def GetTaskStatus(self, request, context):
        loggy.info("stub for future release")

        return oracle_pb2.GetTaskStatusResponse()
