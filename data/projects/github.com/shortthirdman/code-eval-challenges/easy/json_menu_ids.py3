import fileinput, re

json = re.compile('"id": (\d+),')
for line in filter(lambda x: x != '\n', fileinput.input()):
    print(sum(int(i) for i in json.findall(line)))
