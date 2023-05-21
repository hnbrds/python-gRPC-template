import os
import sys
import argparse
from datetime import datetime
import yaml
import logging
from concurrent import futures
import pprint
import grpc
sys.path.append("./proto")
# TODO : import your generated proto module sample_pb2 -> your_service_pb2
from proto import sample_pb2
from proto import sample_pb2_grpc

class Servicer(sample_pb2_grpc.SampleServiceServicer):
    def __init__(self, cfg, logger):
        self.cfg = cfg
        self.logger = logger

    def log_current_time(self):
        req_time = datetime.now().strftime("%Y-%m-%d-%H-%M-%S-%f")
        self.logger.info(f"Client request for {sys._getframe(1).f_code.co_name} at {req_time}")

    # TODO : modify function names - match gRPC proto
    def AddRequest(self, req, ctx):
        self.log_current_time()
        operand_1 = req.operand_1
        operand_2 = req.operand_2
        return sample_pb2.ResultResponse(result=operand_1+operand_2)

    def SubRequest(self, req, ctx):
        self.log_current_time()
        operand_1 = req.operand_1
        operand_2 = req.operand_2
        return sample_pb2.ResultResponse(result=operand_1-operand_2)


def serve(args, cfg, logger):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=cfg['MAX_WORKERS']))
    # TODO : modify sample_pb2_grpc -> your_service_pb2_grpc
    sample_pb2_grpc.add_SampleServiceServicer_to_server(Servicer(cfg, logger), server)
    server.add_insecure_port('[::]:' + str(args.port))
    server.start()
    logger.info("gRPC server initialized")
    server.wait_for_termination()

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--port", type=str, default=13271, help="Server port number")
    parser.add_argument("--config", default="./config.yaml", help="Path to .yaml config file")
    args = parser.parse_args()
    return args

def init_logger():
    logger = logging.getLogger("gRPC SERVER")
    formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
    logger.setLevel(logging.DEBUG)
    stream_handler = logging.StreamHandler()
    stream_handler.setFormatter(formatter)
    logger.addHandler(stream_handler)
    return logger

def load_config(config):
    with open(config) as f:
        cfg = yaml.load(f, Loader=yaml.SafeLoader)
    print("CONFIG INFORMATION")
    pprint.PrettyPrinter(indent=2).pprint(cfg)
    return cfg

if __name__=="__main__":
    args = parse_arguments()
    cfg = load_config(args.config)
    logger = init_logger()
    serve(args, cfg, logger)
