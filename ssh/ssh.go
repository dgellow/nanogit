package ssh

import (
	"github.com/qrclabs/sshgit"

	"github.com/qrclabs/nanogit/config"
)

func checkPubKey(key ssh.PublicKey) (*ssh.PublicKey, error) {
	keystr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))

	filename := "authorized_keys.txt"
	file, err := os.Open(filename)
	if err!= nil {
		log.Fatal(4, "Cannot open file: %s, error: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		fmt.Printf("line: %s\n", line)
		fmt.Printf("keystr: %s\n", keystr)
		fmt.Println("")

		if line == keystr {
			fmt.Println("found key!")
			return &key, nil
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(4, "Error while reading file: %s, error: %v", file, err)
	}

	fmt.Println("found nothing :(")
	return nil, errors.New("key not found")
}

func Listen() {
	sshgit.Listen()
	fmt.Println("Start program")

	config := sshgit.ServerConfig{
		Host: "localhost",
		Port: 1337,
		PrivatekeyPath: "key.rsa",
		KeygenConfig: sshgit.SSHKeygenConfig{"rsa", ""},
		PublicKeyCallback: PublicKeyHandler,
	}

	fmt.Println("Run server")
	sshgit.Listen(config)
}
