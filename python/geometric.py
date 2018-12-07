##### Finding Geometric Sequences

# Given a list (l) and a ratio (r), find the number of groups of 3 indices (i,j,k) in the list such that: 
# 1. i < j < k
# 2. {l[i], l[j], l[k]} is a geometric sequence with a common ratio r
# i.e. r*l[i] == l[j], r*l[j] == l[k] 

#we need a special case for when r is 1 because in this case, the numbers are all the same and so can be in any order. This makes the brute force solution used for other cases too slow.
#when r is 1, we create a list of each unique value, and then we count how many times each value appears and use the formula f(count)=((count)*(count-1)*(count-2))/6 to find the total number of solutions for each unique value, then total them.
#in any other case besides r=1, we simply assume that any index can be i, then we use brute force to search for valid j indicies and then for valid k indicies for each valid j index.

def find_geo_seq(l, r):
    if r==1:
        list=[]
        total = 0
        for x in l:
            go = True
            for y in list:
                if y==x:
                    go=False
            if go==True:
                list.append(x)
                count = 0
                for z in l:
                    if z==x:
                        count+=1
                total = total + ((count*(count-1)*(count-2))/6)
        return total
    else:
        total = 0
        index = 0
        for x in l:
            secondIndex = index+1
            for y in l[secondIndex:]:
                if y==x*r:
                    for z in l[secondIndex+1:]:
                        if z==y*r:
                            total+=1
                secondIndex+=1
            index+=1
        return total

test_cases = [
    ([1, 2, 2, 4], 2, 2),
    ([1, 1, 5, 25, 25, 125, 625], 5, 8),
    ([125, 125, 25, 25, 5], 5, 0),
    ([1, 3, 9, 9, 9, 9, 9, 10, 27, 81], 3, 15),
    ([345]*10000, 1, 166616670000),
    ([1, 1, 1, 1] + [3, 3, 3, 3], 1, 8)
]

for case in test_cases:
    l, r, output = case
    
    print(find_geo_seq(l, r) == output)