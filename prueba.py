import dis

def prueba():
    x = 0
    lista = [0,1,2,3,4,5,6,7,8,9]
    while (x < 10):
        if(x%2==0):
            print(lista[x])
        x = x + 1

dis.dis(prueba)

hola = [0,1,2,3]

if (1 in hola):
    print("si")


