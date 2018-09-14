// Copyright © 2018 Dave Jaccard <dave.jaccard@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	ttt "github.com/ave19/tictacgo/lib"
	"github.com/spf13/cobra"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var defaultAiPlayerNames = []string{
	"Alice",
	"Bob",
}

var defaultAiPlayerStrategies = []ttt.Strategy{
	ttt.NewRememberWinningStrategy(),
	ttt.NewRememberLosingStrategy(),
}

// gameCmd represents the game command
var tttCmd = &cobra.Command{
	Use:   "ttt",
	Short: "Start a game of tic tac toe.",
	Long: `Start a game of tic tac toe, attach a couple of players,
	and see who's the best.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WOULD YOU LIKE TO PLAY A GAME")

		// From options...
		numGames, _ := strconv.Atoi(cmd.Flag("numGames").Value.String())
		numPlayers, _ := strconv.Atoi(cmd.Flag("numPlayers").Value.String())

		if cmd.Flag("verbose").Value.String() == "true" {
			ttt.SetTrace(os.Stdout)
		}

		// Machinery...
		var players []*ttt.Player

		for count := 0; count < (2 - numPlayers); count++ {
			p := ttt.NewPlayer()
			p.SetName(defaultAiPlayerNames[count])
			p.SetNumber(byte(count) + 1)
			p.SetStrategy(defaultAiPlayerStrategies[count])
			players = append(players, &p)
		}

		ttt.Trace.Printf("players -> %#v\n", players)

		for game := 0; game < numGames; game++ {
			b := ttt.NewBoard()
			fmt.Println("\nGame", game+1)
			roundOver := false

			for count := 0; !roundOver; count++ {
				fmt.Println("----------------------------------------")
				fmt.Printf("Round %d:\n", count)

				playerSelector := count % len(players)
				ttt.Trace.Printf("playerSelector = %d\n", playerSelector)

				player := players[playerSelector]
				fmt.Println(" - Current Player:", player.Name())
				ttt.Trace.Println("playerSelector:", playerSelector)
				fmt.Println(b.Int())
				fmt.Println(b.String())
				choice := player.Move(b)
				ttt.Trace.Printf("Choice is %d\n", choice)
				b.Move(choice, player.Number())
				if count >= int(b.NumSquares) {
					roundOver = true
				}
				win, mark := b.Winner()
				if win {
					fmt.Println("==============================================")
					for _, p := range players {
						if p.Number() == mark {
							fmt.Println("Winner: ", p.Name())
							p.Win(b)
						} else {
							fmt.Println("Loser: ", p.Name())
							p.Lose(b)
						}
					}
					roundOver = true
				} else {
					remainingSquares := len(b.ListEmptySquares())
					if remainingSquares == 0 {
						fmt.Println("It was a tie.")
						roundOver = true
						for _, p := range players {
							p.Draw(b)
						}
					}
				}
			}

			fmt.Println("Final Board:")
			fmt.Println(b.Int())
			fmt.Println(b.String())

			for _, p := range players {
				p.Stats(os.Stdout)
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(tttCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tttCmd.Flags().IntP("numGames", "n", 1, "Number of games to play.")
	tttCmd.Flags().IntP("numPlayers", "p", 0, "Number of players.")
	tttCmd.Flags().BoolP("verbose", "v", false, "Verbose logging.")

}
