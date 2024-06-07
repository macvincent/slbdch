# Given dictionaries
m1 = {0:49, 1:400, 2:994, 3:1198, 4:803, 5:1494, 6:607, 7:54, 8:397, 9:4}
m2 = {0:49, 1:399, 2:906, 3:604, 4:417, 5:747, 6:308, 7:28, 8:191, 9:1}
m3 = {2:88, 3:594, 4:385, 5:746, 6:301, 7:27, 8:205, 9:4}

# Initialize dictionaries to store the ratios
m2_m1_ratio = {}
m3_m1_ratio = {}

# Calculate m2/m1 and m3/m1 ratios
for key in m1.keys():
    m2_m1_ratio[key] = m2.get(key, 0) / m1[key] if key in m2 else 0
    m3_m1_ratio[key] = m3.get(key, 0) / m1[key] if key in m3 else 0

# Print the results
print("m1_reqs = ", m1)
print("m2_ratio = ", m2_m1_ratio)
print("m3_ratio = ", m3_m1_ratio)
