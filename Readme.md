![Go](https://img.shields.io/badge/Go-15-blue?logo=go&logoColor=white) 
#  TUI Iptables Manager

Un **TUI** pour gérer facilement les règles de firewall sous **iptables** .  
L’objectif est de rendre la gestion des règles plus lisible et interactive, sans avoir à taper de longues commandes complexes.

---

##  Fonctionnalités

-  Lister les règles existantes par table et chaîne.  
- Naviguer dans les règles avec une interface terminale.  
-  Ajouter de nouvelles règles.  
-  Supprimer une règle sélectionnée.  
-  Sauvegarder ou exporter les règles actuelles.  
-  Colorisation intuitive :  
  - **ACCEPT** = vert  
  - **DROP** = rouge  
  - **REJECT** = jaune  

---

## 🚀 Installation

### Prérequis
- Linux avec **iptables** installé.  
- Accès root (via `sudo`).  
- Un terminal compatible (`xterm`, `alacritty`, `kitty`, …).  

### Build (Go)
```bash
git clone https://github.com/Reentryti/I.T-Manager.git
cd I.T-Manager
go build -o I.T-Manager

