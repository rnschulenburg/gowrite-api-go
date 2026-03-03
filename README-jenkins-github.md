🏆 Einfachste sichere Lösung: GitHub-SSH-Key

Du machst im Grunde:

Jenkins bekommt einen SSH-Schlüssel

GitHub vertraut diesem Schlüssel

Jenkins kann das Repo klonen

✅ Schritt 1 — SSH-Key auf dem Jenkins-Server erzeugen

👉 Auf dem Jenkins-Server einloggen (per SSH)

Dann:

ssh-keygen -t ed25519 -C "jenkins@gowrite-api-go"

Einfach Enter drücken, keine Passphrase.

Ergebnis:

~/.ssh/id_ed25519        ← privater Key (geheim!)
~/.ssh/id_ed25519.pub    ← öffentlicher Key
✅ Schritt 2 — Public Key zu GitHub hinzufügen

Auf dem Jenkins-Server:

cat ~/.ssh/id_ed25519.pub

👉 Kopiere die ganze Zeile.

Jetzt in GitHub:

Repo öffnen

Settings → Deploy Keys

Add deploy key

Key einfügen

„Allow write access“ optional

👉 Fertig ✅

✅ Schritt 3 — Key in Jenkins hinterlegen
Jenkins Weboberfläche:

Manage Jenkins

Credentials

Global → Add Credentials

Typ wählen:

👉 SSH Username with private key

Felder ausfüllen:

Username:

git

Private Key:
👉 „Enter directly“

Dann auf dem Server:

cat ~/.ssh/id_ed25519

👉 gesamten Inhalt einfügen

ID merken

z. B.:

github-ssh
✅ Schritt 4 — Pipeline anpassen

Statt HTTPS-URL:

git branch: 'main',
credentialsId: 'github-ssh',
url: 'git@github.com:rnschulenburg/gowrite-api-go.git'
🚀 Fertig 🎉

Jenkins kann jetzt sicher klonen.
