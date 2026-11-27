#!/usr/bin/env bash
# Builds the starting state described in the question, before any task runs:
#   - groups Developer and Tester
#   - Team A: Developer1, Developer2, Developer3 in the Developer group
#   - Team B: Tester1, Tester2 in the Tester group
#   - the web-based Java project Hdfc_Bank-ClientApp, shared with both groups,
#     containing the zip archive Hdfc-Pack.zip
#
# Linux user names cannot contain spaces, so "Developer 1" becomes Developer1.
# Run as ec2-user; every administrative command goes through sudo.
set -euo pipefail

PROJECT=/opt/Hdfc_Bank-ClientApp

# Groups
sudo groupadd Developer
sudo groupadd Tester

# Team A (developers) and Team B (testers). -m creates the home directory;
# -G adds the user to the team group (each user also gets a private group of
# the same name, which is the Amazon Linux default).
for u in Developer1 Developer2 Developer3; do sudo useradd -m -G Developer "$u"; done
for u in Tester1 Tester2; do sudo useradd -m -G Tester "$u"; done

# Some personal work in Developer3's home, so the removal in task (ii) has
# something to remove.
sudo -u Developer3 mkdir -p /home/Developer3/workspace
sudo -u Developer3 bash -c 'echo "TODO: refactor LoginController" > /home/Developer3/workspace/notes.txt'

# The project: developers can write to it, testers can read it.
sudo mkdir -p "$PROJECT/src/main/java/com/hdfc/clientapp" "$PROJECT/src/main/webapp"
sudo tee "$PROJECT/pom.xml" >/dev/null <<'EOF'
<project>
  <modelVersion>4.0.0</modelVersion>
  <groupId>com.hdfc</groupId>
  <artifactId>Hdfc_Bank-ClientApp</artifactId>
  <version>1.0.0</version>
  <packaging>war</packaging>
</project>
EOF
sudo tee "$PROJECT/src/main/java/com/hdfc/clientapp/LoginController.java" >/dev/null <<'EOF'
package com.hdfc.clientapp;

public class LoginController {
    public boolean login(String user, String password) {
        return user != null && password != null;
    }
}
EOF
echo '<html><body><h1>Client Login</h1></body></html>' | sudo tee "$PROJECT/src/main/webapp/index.html" >/dev/null

# Hdfc-Pack.zip: the release bundle that the task inspects.
PACK=$(mktemp -d)
mkdir -p "$PACK/Hdfc-Pack/config" "$PACK/Hdfc-Pack/sql" "$PACK/Hdfc-Pack/logs"
echo "Hdfc_Bank-ClientApp release bundle" > "$PACK/Hdfc-Pack/README.txt"
printf 'db.host=db.internal\ndb.port=5432\n' > "$PACK/Hdfc-Pack/config/app.properties"
echo "CREATE TABLE clients (id INT PRIMARY KEY, name TEXT);" > "$PACK/Hdfc-Pack/sql/schema.sql"
echo "2026-11-20 10:42:01 DEBUG connection pool started" > "$PACK/Hdfc-Pack/logs/debug.log"
(cd "$PACK" && zip -qr Hdfc-Pack.zip Hdfc-Pack)
sudo mv "$PACK/Hdfc-Pack.zip" "$PROJECT/"
rm -rf "$PACK"

sudo chown -R root:Developer "$PROJECT"
sudo chmod -R g+rwX,o+rX "$PROJECT"
sudo find "$PROJECT" -type d -exec chmod g+s {} +   # new files inherit the Developer group

echo "Setup complete."
