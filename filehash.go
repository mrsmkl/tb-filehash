package main
import "fmt"
import "golang.org/x/crypto/sha3"
import "encoding/hex"
import "os"
// import "io/ioutil"
import "io"
// import "hash"

func makeMerkle(arr [][]byte, idx int, level uint) []byte {
    zeroword := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
    if level == 0 {
        if idx < len(arr) {
            return arr[idx]
        } else {
            return zeroword
        }
    } else {
        hash := sha3.NewLegacyKeccak256()
        hash.Write(makeMerkle(arr, idx, level-1))
        hash.Write(makeMerkle(arr, idx + (1 << (level-1)), level-1))
        return hash.Sum(nil)
    }
}

func getHash(arr [][]byte) []byte {
    return makeMerkle(arr, 0, depth(uint(len(arr)*2 - 1)))
}

func depth(x uint) uint {
    if x <= 1 {
        return 0
    } else {
        return 1 + depth(x/2)
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

    zeroword := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
    
    if len(os.Args) < 2 {
        fmt.Println("give file as an argument")
        return
    }
    fname := os.Args[1]
    fmt.Println("Reading file", fname)
    
    // open file
    f, err := os.Open(fname)
    check(err)

    // read the file
    arr := [][]byte{}
    
    last := 1
    for last != 0 {
        w := make([]byte, 16)
        // n, err := io.ReadAtLeast(f, w, 16)
        n, err := f.Read(w)
        if err == io.EOF {
            last = 0
        } else {
            check(err)
        }
        last = n
        arr = append(arr, w)
    }

    // need to have at least two elements now
    for len(arr) < 2 {
        arr = append(arr, zeroword)
    }

    str := getHash(arr)
    fmt.Println("hello world", hex.EncodeToString(str))
}

