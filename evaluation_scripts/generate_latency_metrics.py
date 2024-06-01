import argparse
import numpy as np

parser = argparse.ArgumentParser("Generate latency metrics")
parser.add_argument("filename", help="name of files containing latencies", type=str)
parser.add_argument('--multiple-units', help="Older version of code with units in file",action=argparse.BooleanOptionalAction)

args = parser.parse_args()

print("Grabbing latency numbers from:", args.filename)
f = open(args.filename, "r")
print(args)
latencies = []
units = set()
for line in f.readlines():
    line = line.strip()
    factor = 1
    if args.multiple_units:
        unit = line[-2:]
        # Put everything into nanoseconds
        if unit == "µs":
            factor = 1000
        elif unit == "ns":
            factor = 1
        else:
            raise Exception("Unit ", unit, " does not have conversion factor yet")
        units.add(unit)
        line = line[:-2]
    latency = float(line) * factor
    latencies += [latency]

print("Units that are present:", units)
latencies = np.array(latencies)

print("99% tail latency is ", np.percentile(latencies, 99), " ns")
print("Average latency is ", np.average(latencies), " ns")
