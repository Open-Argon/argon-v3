class obj1:
    def __init__(self) -> None:
        pass

    def __add__(self, other):
        print("obj1")
        return 10

class obj2:
    def __init__(self) -> None:
        pass

    def __add__(self, other):
        print("obj2")
        return 20

obj1 = obj1()
obj2 = obj2()

print(obj1 + obj2)
print(obj2 + obj1)