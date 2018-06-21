d = {"Michael": 95, "Bob":75, "Tracy": 85}
print(d["Michael"])
print(d["Bob"])
print(d["Tracy"])

d["liwei"] = 28
d["Jack"] = 108
print(d)
print("liwei" in d)
print("liwie" in d)
print(d.get("liwei"))
print(d.get("liwie"), -1)

d.pop("liwei")
print(d)
