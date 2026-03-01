# ShrimPG 🦐 

**S**afe **H**ardened **R**esilient **I**mmutable **M**emory-efficient **P**assword **G**enerator

A minimal, blazingly fast, and slightly opinionated password manager written in Go. It stores your credentials in a local JSON file because the cloud is just someone else's computer (and we don't trust them).

---

## 🛠 Features (or what's not broken yet)
- [x] **Store**: Save passwords before you forget them again.
- [x] **List**: See all your secrets in one place. All your secrets. Safe and easy. 
- [x] **Delete**: Burn the evidence.
- [x] **Generate**: Create passwords so complex even you won't be able to log in.
- [x] **Encryption**: Actually locking the door instead of just closing it.

## 🚀 Installation (The "I'm a Pro" way)

Since you're here, you probably have Go installed. If not... *really?*

```bash
git clone [https://github.com/Krev3tka/ShrimPG.git](https://github.com/Krev3tka/ShrimPG.git)
cd ShrimPG
go run cmd/passwordManager/main.go --help
```

### "I don't want to compile stuff"

Binary releases for Windows/Linux/macOS are coming soon. Hang tight, or just install Go, it's 2026!


## 📖 Usage

Try these before filing a "it doesn't work" issue:

    Add a password: go run cmd/passwordManager/main.go create github my_secret_pass

    List everything: go run cmd/passwordManager/main.go list

    Delete a service: go run cmd/passwordManager/main.go delete twitter (Good riddance!)

    Catch a password: go run cmd/passwordManager/main.go catch reddit (important thing!)
    
## ⚠️ Troubleshooting / FAQ

Q: Where is my passwords.json? A: It's in the project folder. If it's missing, you probably deleted it. Nice job.

Q: Is it secure? A: Currently, it's a plain JSON. If someone steals your laptop, they own your life. Encryption is coming, but for now, maybe don't leave your PC at cafe?

Q: Why "ShrimPG"? A: Because shrimps are cool, and "ShrimPM" sounded like a late-night medication.


## 🛡 Security Note

We've added passwords.json to .gitignore so you don't accidentally leak your "123456" password to the entire world. You're welcome.
