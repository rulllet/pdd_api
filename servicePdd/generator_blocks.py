def block_list(number: int) -> list:
    return [n for s in
            [[i for i in range(w, w + 5)]
             for w in range(number, 801, 20)]
            for n in s]

def format_go(blocks: list) -> str:
    result = "package servicePdd\n\n"
    for _, i in enumerate(blocks):
        result += f"var block{_+1} = [{len(i)}]int64{{{str(i)[1:-1]}}}\n"
    return result

BLOCK1 = block_list(1)
BLOCK2 = block_list(6)
BLOCK3 = block_list(11)
BLOCK4 = block_list(16)

my_file = open("blocks.go", "w+")
my_file.write(format_go([BLOCK1, BLOCK2, BLOCK3, BLOCK4]))
my_file.close()