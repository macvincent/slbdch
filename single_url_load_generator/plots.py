import matplotlib.pyplot as plt

# Sample dictionaries
n_requests = {0:49, 1:400, 2:900, 3:1198, 4:800, 5:1495, 6:604, 7:302, 8:51, 9:397, 10:13, 11:1}
k000 =   {0:49, 1:400, 2:899, 3:1125, 4:786, 5:1263, 6:585, 7:302, 8:52, 9:397, 10:13, 11:1}
k025 = {0:49, 1:400, 2:898, 3:979, 4:693, 5:1137, 6:556, 7:301, 8:52, 9:398, 10:12, 11:1}
k050 = {0:49, 1:399, 2:899, 3:785, 4:468, 5:809, 6:294, 7:250, 8:42, 9:397, 10:13, 11:1}
k065 = {0:49, 1:400, 2:898, 3:666, 4:396, 5:750, 6:304, 7:147, 8:30, 9:277, 10:13, 11:1}
k075 = {0:49, 1:400, 2:835, 3:597, 4:392, 5:782, 6:302, 7:143, 8:32, 9:201, 10:4}
k090 =  {0:49, 1:399, 2:765, 3:607, 4:419, 5:765, 6:296, 7:153, 8:32, 9:197, 10:3}

# Function to create a plot for a given dictionary with a universal max key
def plot_dict(data, max_key):
    x_values = list(range(max_key + 1))
    y_values = [data.get(x, 0) for x in x_values]
    return x_values, y_values

# Calculate the universal max key
max_key = max(n_requests.keys())

# Plot d1
x1, y1 = plot_dict(n_requests, max_key)
plt.figure(figsize=(10, 5))
plt.plot(x1, y1, marker='o')
plt.axhline(y=1000, color='r', linestyle='--', label='Threshold')  # Add horizontal line at y=1000
plt.title('Number of requests for each one second bucket')
plt.xlabel('Timestamp bucket (second)')
plt.ylabel('Number of Requests')
plt.grid(True)
plt.xticks(range(max_key + 1))  # Ensure all integers are shown on the x-axis
plt.legend()
plt.show()

# Plot the dictionaries for different values of k

plt.figure(figsize=(10, 5))

x000, y000 = plot_dict(k000, max_key)
plt.plot(x000, y000, marker='o', label='γ=0')

x025, y025 = plot_dict(k025, max_key)
plt.plot(x025, y025, marker='^', label='γ=0.25')

x050, y050 = plot_dict(k050, max_key)
plt.plot(x050, y050, marker='d', label='γ=0.50')

x065, y065 = plot_dict(k065, max_key)
plt.plot(x065, y065, marker='p', label='γ=0.65')

x075, y075 = plot_dict(k075, max_key)
plt.plot(x075, y075, marker='p', label='γ=0.75')

x090, y090 = plot_dict(k090, max_key)
plt.plot(x090, y090, marker='h', label='γ=0.90')

plt.title('Number of Requests Served by Server 1 in each one second bucket')
plt.xlabel('Timestamp bucket (second)')
plt.ylabel('Number of Requests Served by Server 1')
plt.grid(True)
plt.xticks(range(max_key + 1))  # Ensure all integers are shown on the x-axis
plt.legend()
plt.show()
