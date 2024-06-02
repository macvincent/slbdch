from matplotlib import pyplot as plt
import numpy as np

# Load csv data from file as numpy array
# plt.plot(0, 0, 'bo')
experiments = [1, 10, 100, 1000, 10000, 100000, 1000000]
value = []
for replica_count in experiments:
    # file_name = f'./kademlia_vs_cycle/kademlia/1-{replica_count}.txt'
    file_name = f'./kademlia_vs_cycle/kademlia/latency_{replica_count}_replicas.csv'
    data = np.loadtxt(file_name, delimiter=',')
    # sort data in descending order
    data = np.sort(data)
    value.append(np.percentile(data, 90))
    plt.plot(np.log10(replica_count), value[-1], 'ro')
     
plt.plot(np.log10(experiments), value, 'r-')

# plt.plot(0, 0, 'bo')
value = []
for replica_count in [1, 10, 100, 1000, 10000, 100000, 1000000]:
    file_name = f'./kademlia_vs_cycle/cycle/latency_{replica_count}_replicas.csv'
    data = np.loadtxt(file_name, delimiter=',')
    # blue color
    data = np.sort(data)
    value.append(np.percentile(data, 90))
    plt.plot(np.log10(replica_count), value[-1], 'bo')
plt.plot(np.log10(experiments), value, 'b-')

plt.show()

    
