def interpret(code):
    memory = {}
    pointer = 0
    code_ptr = 0
    loops = []

    while code_ptr < len(code):
        command = code[code_ptr]

        if command == '>':
            pointer += 1
        elif command == '<':
            pointer -= 1
        elif command == '+':
            if pointer not in memory:
                memory[pointer] = 0
            memory[pointer] += 1
        elif command == '-':
            if pointer not in memory:
                memory[pointer] = 0
            memory[pointer] -= 1
        elif command == '.':
            print(chr(memory.get(pointer, 0)), end='')
        elif command == ',':
            memory[pointer] = ord(input())
        elif command == '[':
            if memory.get(pointer, 0) == 0:
                loop_depth = 1
                while loop_depth > 0:
                    code_ptr += 1
                    if code[code_ptr] == '[':
                        loop_depth += 1
                    elif code[code_ptr] == ']':
                        loop_depth -= 1
            else:
                loops.append(code_ptr)
        elif command == ']':
            if memory.get(pointer, 0) != 0:
                code_ptr = loops[-1]
            else:
                loops.pop()
        code_ptr += 1

# Example usage
interpret("++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.")
# Output: Hello World!
