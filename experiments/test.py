from matplotlib import pyplot as plt
import numpy as np
nodes = [1, 3, 5]
node_to_color = {1: 'r', 3: 'g', 5: 'b'}
csv_file = './load_tester_throughput.csv'

results_to_be_plotted = []
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
    
    runs = [100, 315, 1000, 3160, 10000]
    results = []
    for run in runs:
        throughput = runs_to_throughput[run]
        average_throughput = np.array(throughput).max()
        results.append(average_throughput)
        print(f'Nodes {node}, Run {run}, Max Throughput: {average_throughput}')
        results_to_be_plotted.append((node, run, average_throughput))
        plt.plot(run, average_throughput, f'{node_to_color[node]}o')
    plt.plot(runs, results, f'{node_to_color[node]}-', label=f'{node} Workers')
    plt.yscale('log')
    plt.xscale('log')
    plt.xlabel('Request Throughput')
    plt.ylabel('Response Throughput')
    plt.legend()
    
for node in nodes:
    plot_nodes(node)
plt.show()

with open(csv_file, 'w') as f:
    for node, run, throughput in results_to_be_plotted:
        f.write(f'{node},{run},{throughput}\n')