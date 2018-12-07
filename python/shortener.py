######## Shortener

# Given two strings s, t, determine if string s can be changed into string t using the following rule:

# Remove any number of lowercase letters from s (or none at all), then capitalize all remaining letters of s.

# If it is possible to create the string t by processing s using the rule, then the function should return True, otherwise return False.

#my general strategy was to compare the first character of s and t. If they are the same letter, then we drop both and move on to the next character of s and of t.
#If they are not the same, we see if the character from s can be dropped. If not, return false. If we can, we do that and move on to the next character from s.
#we do this until s and t are equal length. Then we simply compare s.upper() to t, and return true if they are the same and false if not.

def shortener(s, t):
    while len(s)>len(t):
        if s[0:1].upper()==t[0:1]:
            t=t[1:]
        elif not s[0:1].islower():
            return False
        s=s[1:]
    return True if s.upper()==t else False

# Test Cases
test_cases = [
    ("daBccd", "ABC", True),
    ("sYOCa", "YOCN", False),
    ("aaaaaa", "AAAAAAA", False),
    ("SVAHHHMVIIDYIcOSHMDUAVJRIBxBZQSUBIVEBHfVTZVSHATUYDJGDRRUBQFHEEEUZLQGXTNKFWUYBAeFKUHSFLZEUINBZYRIXOPYYXAEZZWELUPIEIWGZHEIYIROLQLAVHhMKRDSOQTJYYLTCTSIXIDAnPIHNXENWFFZFJASRZRDAPVYPAViVBLVGRHObnwlcyprcfhdpfjkyvgyzpovsgvlqbhtwrucvszaqinbgeafuswkjrcexvyzq","SVAHHHMVIIDYIOSHMDUAVJRIBBZQSUBIVEBHVTZVSHATUYDJGDRRUBQFHEEEUZLQGXTNKFWUYBAFKUHSFLZEUINBZYRIXOPYYXAEZZWELUPIEIWGZHEIYIROLQLAVHMKRDSOQTJYYLTCTSIXIDAPIHNXENWFFZFJASRZRDAPVYPAVVBLVGRHO", True),
    ("a", "AA", False),
    ("UZJMUCYHpfeoqrqeodznwkxfqvzktyomkrVyzgtorqefcmffauqhufkpptaupcpxguscmsbvolhorxnjrheqhxlgukjmgncwyastmtgnwhrvvfgbhybeicaudklkyrwvghpxbtpyqioouttqqrdhbinvbywkjwjkdiynvultxxxmwxztglbqitxmcgiusfewmsvxchkryzxipbmgrnqhfmlghomfbsKjglimxuobomfwutwfcmklzcphbbfohnaxgbaqbgocghaaizyhlctupndmlhwwlxxvighhjjrctcjBvxtagxbhrbrWwsyiiyebdgyfrlztoycxpjcvmzdvfeYqaxitkfkkxwybydcwsbdiovrqwkwzbgammwslwmdesygopzndedsbdixvi","UZJMUCYH", False)
]

for case in test_cases:
    s, t, output = case
    print(shortener(s, t) == output)

