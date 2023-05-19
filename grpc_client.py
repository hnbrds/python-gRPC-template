import os
import sys
import grpc
sys.path.append("./proto")
import yaml
import argparse

# TODO : import your own generated grpc modules
from proto import sample_pb2
from proto import sample_pb2_grpc
from google.protobuf.json_format import MessageToDict

def main(args):
    addr = args.ip + ":" + str(args.port)
    channel = grpc.insecure_channel(addr)
    # TODO : make stub and call your own declared function
    stub = sample_pb2_grpc.SampleServiceStub(channel)
    req = sample_pb2.OperandMessage(
        operand_1 = 3,
        operand_2 = 5
    )
    res = stub.AddRequest(req)
    print(res)

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--ip", default="127.0.0.1", help="IP address")
    parser.add_argument("--port", default=13271, help="Port number")
    args = parser.parse_args()
    return args

if __name__=="__main__":
    args = parse_arguments()
    main(args)

