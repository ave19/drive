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
		players := []*ttt.Player{}
		//games := []*ttt.Game{}
		var gameChannels []chan ttt.Game
		ttt.Trace.Println("GameChannels", gameChannels)

		// Create machine players...
		// A "human" would use a specical strategy, which we don't have yet...
		// setting this to zero makes two ai players
		//numPlayers = 0

		for i := numPlayers; i < 2; i++ {
			ttt.Trace.Println("Making player:", i+1)
			p := ttt.NewPlayer()
			p.SetName(defaultAiPlayerNames[i])
			//p.SetNumber(byte(i + 1))
			p.SetStrategy(defaultAiPlayerStrategies[i])
			//gameChannel := make(chan ttt.Game)
			//gameChannels = append(gameChannels, gameChannel)
			//ttt.Trace.Println("gameChannels", gameChannels)
			go p.Play()
			<-p.ReadyChan
			players = append(players, &p)
		}

		ttt.Trace.Println("len players:", len(players))
		for i, p := range players {
			fmt.Printf("Player %d:\n", i+1)
			fmt.Println("  Name:", p.Name())
			fmt.Println("  Strategy:", p.Strategy().Name())
		}

		for num := 0; num < numGames; num++ {

			var game ttt.Game
			// Randomly pick which player gets to go first
			if (rand.Int31n(2)) == 1 {
				game = ttt.NewGame(*players[0], *players[1])
			} else {
				game = ttt.NewGame(*players[1], *players[0])
			}

			game.Start()

			/*
				b := ttt.NewBoard()
				fmt.Println("\nGame", game+1)
				roundOver := false
				ttt.Trace.Println("bChannels", bChannels)
				for count := 0; !roundOver; count++ {
					fmt.Println("----------------------------------------")
					fmt.Printf("Round %d:\n", count)
					playerSelector := count % len(players)
					player := players[playerSelector]
					fmt.Println(" - Current Player:", player.Name())
					ttt.Trace.Println("playerSelector:", playerSelector)
					fmt.Println(b.Int())
					fmt.Println(b.String())
					ttt.Trace.Println("Sending", b.Int(), "to channel", bChannels[playerSelector])
					bChannels[playerSelector] <- b
					ttt.Trace.Println("Listening to moveChan", player.MoveChan)
					b.Move(<-player.MoveChan, player.Number())
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
			*/
		}

		for _, p := range players {
			p.Stats(os.Stdout)
			fmt.Println()
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
