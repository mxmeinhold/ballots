package controllers

import (
    "github.com/revel/revel"
    "net/http"
    "io/ioutil"
    "io"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "encoding/csv"
    "bufio"
    "os"
)

type App struct {
    *revel.Controller
}

type PullRequest struct {
    Title string
    Number int
    Html_url string
    User User `json:"user"`
    Body string
}

type User struct {
    Login string `json:"login"`
    Html_url string `json:"html_url"`
    Avatar_url string `json:"avatar_url"`
}

func (c App) Index() revel.Result {
    // Fetch open pull requests
    resp, err := http.Get("https://api.github.com/repos/ComputerScienceHouse/Constitution/pulls")
    if err != nil {
        fmt.Printf("Error fetching github information")
        return c.Render()
    }

    // Read http response body
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response body")
        return c.Render()
    }

    // Cast response body to byte array for json parsing
    respString := []byte(string(body))
    // Create empty array of pull request structs
    var prs []PullRequest

    // Parse json automagically into pull request array (pointers are neat)
    err = json.Unmarshal(respString, &prs)
    if err != nil {
        fmt.Printf("Error parsing json")
        return c.Render()
    }

    return c.Render(prs)
}

func (c App) CustomBallots(prompt string, answers string) revel.Result {
    // TODO
    // Variable number of pokemon/ballots
    pokemons := getPokemon(100)
    // Delimit text input based on new lines in the text area
    options := strings.Split(answers, "\n")
    return c.Render(prompt, options, pokemons)
}

func (c App) Ballots(prnumber int, numballots int) revel.Result {
    // Fetch pull request diff
    resp, err := http.Get("https://patch-diff.githubusercontent.com/raw/ComputerScienceHouse/Constitution/pull/" +
        strconv.Itoa(prnumber) + ".diff")
    if err != nil {
        fmt.Printf("Error fetching PR diff")
        return c.Render()
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response body")
        return c.Render()
    }

    diffString := string(body)
    // Fix new lines
    strings.Replace(diffString, `\n`, "\n", -1)

    // Fetch PR data to determine the title
    resp, err = http.Get("https://api.github.com/repos/ComputerScienceHouse/Constitution/pulls/" + strconv.Itoa(prnumber) + ".diff")
    if err != nil {
        fmt.Printf("Error fetching PR title")
        return c.Render()
    }

    defer resp.Body.Close()
    titleBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response body")
        return c.Render()
    }

    // Cast response body to usable byte array
    titleString := []byte(string(titleBody))
    var pr PullRequest
    // Parse title with json
    err = json.Unmarshal(titleString, &pr)
    if err != nil {
        fmt.Printf("Error parsing json")
        return c.Render()
    }

    // TODO
    // Variable number of pokemon/ballots
    pokemons := getPokemon(100)
    return c.Render(diffString, pokemons, pr)
}

func getPokemon(numballots int) []string {
    // Open pokemon csv file
    pokefile, err := os.Open(os.Getenv("PCSV_PATH"))
    if err != nil {
        fmt.Printf("Error opening pokemon.csv")
        return nil
    }

    r := csv.NewReader(bufio.NewReader(pokefile))
    // Off by one as usual
    numballots = numballots + 1
    // Create empty pokemon array (this is some weird go call that just works)
    pokemons := make([]string, numballots)
    for i := 1; i < numballots; i++{
        pokemon, err := r.Read()
        if err == io.EOF {
            break
        }
        pokemons[i] = pokemon[1]
    }
    // Closing file is important
    pokefile.Close()
    return pokemons
}
