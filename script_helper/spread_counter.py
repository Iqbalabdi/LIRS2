import sys

if len(sys.argv) < 2:
    print("Usage: python spread_counter.py <filename>")
    sys.exit(1)

LIMIT = 100

count = {}
spread = {LIMIT+1: 0}
lines = 0

filepath = sys.argv[1]
filename = filepath.split("/")[-1]
output = open(f'spread_count_{filename}', 'w')

source = open(filepath, 'r')
for line in source.readlines():
    lba = line.split(",")[0]
    if lba in count:
        count[lba] += 1
    else:
        count[lba] = 1
    lines += 1

for lba in count:
    if count[lba] > LIMIT:
        spread[LIMIT+1] += 1
        continue
    if count[lba] in spread:
        spread[count[lba]] += 1
    else:
        spread[count[lba]] = 1
    
# group_by = 5
# group_count = 0

# for i in range(1, LIMIT+1):
#     if i in spread:
#         group_count += spread[i]
#     else:
#         group_count += 0
#     if i % group_by == 0:
#         output.write(f'{int(i-group_by+1)}-{i},{group_count}\n')
#         group_count = 0
# else:
#     output.write(f'>{LIMIT},{spread[LIMIT+1]}\n')

for i in range(1, LIMIT+1):
    if i in spread:
        output.write(f'{i},{spread[i]}\n')
    else:
        output.write(f'{i},0\n')
else:
    output.write(f'>{LIMIT},{spread[LIMIT+1]}\n')

# print(lines)
# print(spread)