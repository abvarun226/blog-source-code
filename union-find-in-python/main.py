class MergeAccounts:
    def __init__(self, accounts):
        self.id = {}
        self.owner = {}
        self.size = {}
        self.accounts = accounts

        for account in accounts:
            name = account[0]
            for email_id in account[1:]:
                self.owner[email_id] = name
                self.id[email_id] = email_id
                self.size[email_id] = 1

    def merge_accounts(self):
        for account in self.accounts:
            p = account[1]
            for q in account[2:]:
                self.Union(p, q)

        temp = {}
        for k in self.id:
            v = self.Find(k)
            temp.setdefault(v, []).append(k)

        merged = []
        for k, v in temp.items():
            name = self.owner[k]
            merged.append([name] + sorted(v))

        return merged

    def Find(self, i):
        while i != self.id[i]:
            # optimization: path compression, i.e. point node to it's grandparent.
            # keeps the tree completely flat.
            self.id[i] = self.id[self.id[i]]
            i = self.id[i]
        return i

    def Union(self, p, q):
        i = self.Find(p)
        j = self.Find(q)
        if i == j:
            return
        # optimization: always put the smaller tree lower, i.e. root of smaller tree is bigger tree.
        if self.size[i] < self.size[j]:
            self.id[i] = j
            self.size[j] += self.size[i]
        else:
            self.id[j] = i
            self.size[i] += self.size[j]


if __name__ == '__main__':
    accounts = [
        ["John", "johnsmith@mail.com", "john_newyork@mail.com"],
        ["John", "johnsmith@mail.com", "john00@mail.com"],
        ["Mary", "mary@mail.com"],
        ["John", "johnnybravo@mail.com"],
    ]

    m = MergeAccounts(accounts)
    print(m.merge_accounts())

    accounts = [
        ["Gabe", "Gabe0@m.co", "Gabe3@m.co", "Gabe1@m.co"],
        ["Kevin", "Kevin3@m.co", "Kevin5@m.co", "Kevin0@m.co"],
        ["Ethan", "Ethan5@m.co", "Ethan4@m.co", "Ethan0@m.co"],
        ["Hanzo", "Hanzo3@m.co", "Hanzo1@m.co", "Hanzo0@m.co"],
        ["Fern", "Fern5@m.co", "Fern1@m.co", "Fern0@m.co"],
    ]
    m = MergeAccounts(accounts)
    print(m.merge_accounts())

    accounts = [
        ["Alice", "alice@example.com", "alice_wonderland@example.com"],
        ["Bob", "bob@example.com", "bob_marley@example.com"],
        ["Alice", "alice@mail.com", "alice@example.com"],
        ["Alice", "alice@mail.com"],
    ]
    m = MergeAccounts(accounts)
    print(m.merge_accounts())
