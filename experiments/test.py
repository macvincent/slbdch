from matplotlib import pyplot as plt
import numpy as np
nodes = [1, 3, 5]
csv_file = './load_tester_throughput.csv'

def plot_nodes(node):
    file_name = f'./load_tester_throughput_{node}_node.txt'
    runs_to_throughput = {}
    with open(file_name, 'r') as f:
        for line in f:
            run, throughput = line.split()
            run = int(run)
            if run not in runs_to_throughput:
                runs_to_throughput[run] = []
            runs_to_throughput[run].append(float(throughput))
    
    runs = [100, 1000, 10000]
    results = []
    for run in runs:
        throughput = runs_to_throughput[run]
        average_throughput = np.array(throughput).max()
        results.append(average_throughput)
        print(f'Nodes {node}, Run {run}, Max Throughput: {average_throughput}')
        plt.plot(np.log10(run), average_throughput, 'ro')
    plt.plot(np.log10(runs), results, 'b-')
    plt.xlabel('Request Throughput')
    plt.ylabel('Response Throughput')
    
for node in nodes:
    plot_nodes(node)
plt.show()