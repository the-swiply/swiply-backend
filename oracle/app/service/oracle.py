from app.api.oracle import oracle_pb2_grpc, oracle_pb2
from app.houston import loggy
import pandas as pd
import lightfm
import numpy as np
from lightfm.data import Dataset
from lightfm.cross_validation import random_train_test_split
from lightfm import LightFM
from lightfm.evaluation import precision_at_k


class OracleService(oracle_pb2_grpc.OracleServicer):
    def __init__(self) -> None:
        pass

    def RetrainLFMv1(self, request, context):
        loggy.info("start retrain LFMv1")

        return oracle_pb2.RetrainLFMv1Response()

    def GetTaskStatus(self, request, context):
        loggy.info("stub for future release")

        return oracle_pb2.GetTaskStatusResponse()
