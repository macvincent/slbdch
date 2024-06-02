from matplotlib import pyplot as plt
import numpy as np

# Load csv data from file as numpy array
for replica_count in [1, 10, 100, 1000, 10000, 100000, 1000000]:
    file_name = f'./kademlia_vs_cycle/kademlia/latency_{replica_count}_replicas.csv'
    data = np.loadtxt(file_name, delimiter=',')
    plt.plot(np.log(replica_count), np.percentile(data, 99), 'ro')
    
for replica_count in [1, 10, 100, 1000, 10000, 100000, 1000000]:
    file_name = f'./kademlia_vs_cycle/cycle/latency_{replica_count}_replicas.csv'
    data = np.loadtxt(file_name, delimiter=',')
    # blue color
    plt.plot(np.log(replica_count), np.percentile(data, 99), 'bo')

plt.show()

    
