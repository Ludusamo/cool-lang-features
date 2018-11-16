from random import randint, sample
import json

# Read (All), Read (Specific), Write, Update, Delete
percentages=[0.3, 0.3, 0.2, 0.1, 0.1]

registered=set()
usedNames=set()
removed=set()
idToName={}
attackNum=0

def genReadAllAttack():
    return "GET http://localhost:8080/api/feature"

def selectRandomID():
    if len(registered) == len(removed): raise Exception('cannot perform attack')
    ID=sample(registered, 1)[0]
    while ID in removed:
        ID=sample(registered, 1)[0]
    return ID

def genReadSpecificAttack():
    ID=selectRandomID()
    return "GET http://localhost:8080/api/feature/%d" % ID

def genWriteAttack():
    name=str(randint(0, 100000))
    while name in usedNames:
        name=randint(0, 100000)
    usedNames.add(name)
    description=str(randint(0, 100000))
    ID=len(registered) + 2
    registered.add(ID)
    idToName[ID]=name

    data={"Name": name, "Description": description}
    with open('tmp/%d.json' % attackNum, 'w+') as jsonFile:
        json.dump(data, jsonFile)
    return "POST http://localhost:8080/api/feature\n@tmp/%d.json" % attackNum

def genUpdateAttack():
    ID=selectRandomID()

    data={"Name": idToName[ID], "Description": str(randint(0, 100000))}
    with open('tmp/%d.json' % attackNum, 'w+') as jsonFile:
        json.dump(data, jsonFile)
    return "PATCH http://localhost:8080/api/feature/%d\n@tmp/%d.json" % (ID, attackNum)

def genDeleteAttack():
    ID=selectRandomID()
    removed.add(ID)
    del idToName[ID]

    return "DELETE http://localhost:8080/api/feature/%d" % ID

attackGens=[genReadAllAttack, genReadSpecificAttack, genWriteAttack, genUpdateAttack, genDeleteAttack]

def selectAttack():
    rand=randint(0, 100)
    threshold=1
    atk=0
    for i in range(len(percentages)):
        threshold+=100*percentages[i]
        if rand < threshold: return i

with open('tmp/vegeta_attack', 'w+') as f:
    for i in range(30000):
        atk=attackGens[selectAttack()]
        atkStr=""
        while True:
            try:
                atkStr=atk()
                break
            except Exception as e:
                atk=attackGens[selectAttack()]
                pass
        f.write("%s\n\n" % atkStr)
        attackNum+=1
