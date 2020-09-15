import sys
import pymongo
from pymongo import collection
MONGO_URI = 'mongodb+srv://cmater:dvlp@cluster0.j4fgv.gcp.mongodb.net/users?retryWrites=true&w=majority'



class Project():

    def __init__(self):
        self.name = ""
        self.img =""

    def update(self, **kwargs):
        self.name = kwargs[name] if name in kwargs else ""
        self.img = kwargs[img] if img in kwargs else ""


class Person():

    def __init__(self):
        self._id =""
        self.name=""
        self.email=""
        self.education=[]
        self.username=""
        self.password=""
        self.route=""
        self.specialization=[]
        self.projects =[]
        self.achievements= []

    def get(self, key , **kwargs):
        return kwargs[key] if key in kwargs else ""
    

    def update(self , **kwargs ):
        self.name = self.get( 'name' , **kwargs)
        self.email = self.get( 'email' , **kwargs)
        self.username = self.get( 'username' , **kwargs)
        self.password = self.get( 'password' , **kwargs)
        self.route = self.get( 'route' , **kwargs)
        self.education = self.get( 'education' , **kwargs)
        self.specialization = self.get( 'specialization' , **kwargs)
        self.achievements = self.get( 'achievements' , **kwargs)
        self.projects = self.get( 'projects' , **kwargs ) 



        





client = pymongo.MongoClient(MONGO_URI)

db = client['users']
coll = db['info']

person = Person()
person.update(
    name= 'Shuvayan',
    email = 'sgd030@gmail.com',
    username= 'nibba',
    password = 'pass123' ,
    route = 'shuvayan',
    education = ['HVM' , 'JU'],
    specialization = ['CV' ,'NLP'],
    achievements = ['ache1' , 'aceh2'],
    projects = [{'name':'pro1' , 'img':'im1'} , {'name':'pro2' , 'img':'im2' } , {'name':'pro3'}]
)

doc ={}
for k , v in vars(person).items():
    doc[k]=v
del doc['_id']


# write password in hash.txt

with open('hash.txt' , 'w') as f:
    f.write(doc['password'])
f.close()

#write Password at file "hash.txt"

#Hash calculation done
from ctypes import *

hash = cdll.LoadLibrary('./hash.so')
hash.GETHASH()

#Read Hash from "hash.txt"
with open('hash.txt', 'r') as f:
    hash = f.read().strip()
f.close()


print('famous hash : ', hash)

doc['password']=hash






if coll.find_one({"username":doc['username']}) is not None:
    sys.exit()


idd = coll.insert_one(doc).inserted_id
print(idd)

# ob = coll.find_one( {"name":"Shuvayan", "email" : "sgd030@gmail.com"})
# print(ob)







