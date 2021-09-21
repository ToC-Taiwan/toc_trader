'''GO MODULE MASTER'''

with open('./go.mod', "rt") as fin:
    with open('./temp_go.mod', "wt") as fout:
        total = fin.readlines()
        for line in total:
            start = line.find(' v')
            if start == -1:
                fout.write(line)
            else:
                fout.write(line[:start]+' master\n')
