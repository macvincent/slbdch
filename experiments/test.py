from matplotlib import pyplot as plt
import numpy as np

# Load csv data from file as numpy array
plt.plot(0, 0, 'bo')
# for replica_count in [1, 10, 100, 1000, 10000, 100000, 1000000]:
#     file_name = f'./kademlia_vs_cycle/kademlia/1-{replica_count}.txt'
#     data = np.loadtxt(file_name, delimiter=',')
#     plt.plot(np.log(replica_count), np.mean(data), 'ro')
    
# plt.plot(0, 0, 'bo')
for replica_count in [1, 10, 100, 1000, 10000, 100000, 1000000]:
    file_name = f'./kademlia_vs_cycle/cycle/latency_{replica_count}_replicas.csv'
    data = np.loadtxt(file_name, delimiter=',')
    # blue color
    plt.plot(np.log(replica_count), np.mean(data), 'bo')

plt.show()

    
