x=list(map(int,input().split()))
n=x[0]
const=int(2e9)+1
h=[-const]+x[1:]+[-const]
prav=[0]*(n+2)
lev=[0]*(n+2)
st1=[0]
st2=[0]
for i in range(1,n+2):
    while h[st1[-1]]>h[i]:
        if i==n+1:
            prav[st1.pop()]=n+1
        else:
            prav[st1.pop()]=i
    st1.append(i)
    while h[st2[-1]]>h[-(i+1)]:
        if n-i+1!=0:
            lev[st2.pop()]=n-i+1
        else:
            lev[st2.pop()]=0
    st2.append(-(i+1))
bestS=0
for i in range(1,len(h)-1):
    if h[i]*(prav[i]-lev[i]-1)>bestS:
        bestS=h[i]*(prav[i]-lev[i]-1)
print(bestS)



'''st = input()
class Stack:
    def __init__(self):
        self.items = []

    def isEmpty(self):
        return self.items == []

    def push(self, item):
        self.items.append(item)

    def pop(self):
        return self.items.pop()

    def peek(self):
        return self.items[len(self.items)-1]

    def size(self):
        return len(self.items)

    def clear(self):
        self.items = []
        print("ok")

main_stack = Stack()
for i in range(len(st)):
    if st[i] in "0123456789":
        main_stack.push(st[i])'''
    



##class Queue:
##    def __init__(self):
##        self.items = []
##
##    def is_empty(self):
##        return not bool(self.items)
##
##    def push(self, item):
##        self.items.append(item)
##
##    def pop(self):
##        return self.items.pop(0)
##
##    def front(self):
##        return self.items[0]
##
##    def size(self):
##        return len(self.items)
##
##    def clear(self):
##        self.items = []
##        return "ok"
##queue = Queue()
##
##while True:
##    command = input().split()
##
##    if len(command) == 2:
##        cmd, num = command[0], command[1]
##    else:
##        cmd = command[0]
##
##    if cmd == "push":
##        queue.push(int(num))
##        print("ok")
##
##    elif cmd == "pop":
##        print(queue.pop())
##
##    elif cmd == "front":
##        print(queue.front())
##
##    elif cmd == "size":
##        print(queue.size())
##
##    elif cmd == "clear":
##        queue.clear()
##        print("ok")
##
##    elif cmd == "exit":
##        print("bye")
##        break



##class Stack:
##    def __init__(self):
##        self.items = []
##
##    def isEmpty(self):
##        return self.items == []
##
##    def push(self, item):
##        self.items.append(item)
##
##    def pop(self):
##        return self.items.pop()
##
##    def peek(self):
##        return self.items[len(self.items)-1]
##
##    def size(self):
##        return len(self.items)
##
##    def clear(self):
##        self.items = []
##        print("ok")
##
##stack = Stack()
##
##while True:
##    command = input().split()
##    if command[0] == 'push':
##        stack.push(int(command[1]))
##        print("ok")
##    elif command[0] == 'pop':
##        print(stack.pop())
##    elif command[0] == 'back':
##        print(stack.peek())
##    elif command[0] == 'size':
##        print(stack.size())
##    elif command[0] == 'clear':
##        stack.clear()
##    elif command[0] == 'exit':
##        print("bye")
##        break
