package scripts

import (
	"fmt"
	"os"
	"smart-fridge/openai"
	"sync"
)

var greetings = []string{
	"Tür auf, gute Laune rein – was darf’s sein?",
	"Willkommen zurück in der kühlen Komfortzone.",
	"Lass mich raten… was Kaltes gegen innere Leere?",
	"Und, was klaust du mir heute wieder?",
	"Oh wow, Besuch! Ich fühl mich ganz geehrt.",
	"Jemand auf der Suche nach verborgenen Kalorien?",
	"Tür auf, Herz auf – was brauchst du, Freund der Nacht?",
	"Suchst du was oder tust du nur so?",
	"Tadaaa – Bühne frei für deinen nächsten Snack-Fail!",
	"Na los, überrasche mich.",
	"Jemand hat Hunger… oder Langeweile?",
	"Was auch immer du suchst – es ist wahrscheinlich hinter der Milch.",
	"Ach guck, der Türöffner höchstpersönlich!",
	"Achtung, gleich wird’s wieder ungesund… oder?",
	"Hoppala, Besuch aus der Außenwelt!",
	"Du schon wieder! Ich hätte fast Inventur gemacht.",
	"Bereit, mich zu enttäuschen? Ich bin’s auch.",
	"Riechst du das? Das ist der Duft deiner schlechten Entscheidungen.",
	"Willkommen in der Kältekammer des Genusses – oder der Reue.",
	"Hereinspaziert in die Snack-Zone!",
	"Gönnung incoming… ich seh’s schon.",
	"Aha! Der Magen hat wieder das Kommando übernommen.",
	"Wer was will, muss kühl bleiben – du weißt Bescheid.",
	"Wenn du was Gesundes suchst: unten links. Aber das tust du eh nicht, oder?",
	"Na du, auf Schatzsuche im Kühlparadies?",
	"Ayo, was geht bei dir?",
	"Na, wieder auf geheimer Mission unterwegs?",
	"Oha, und da ist er wieder – wie ein Boomerang!",
	"Ich hab dich kommen hören…",
	"Uff, du siehst aus, als bräuchtest du 'ne Umarmung… oder 'nen Snack.",
	"Und, heute schon was Verrücktes gemacht? Außer mich geweckt?",
	"Surprise, ich bin immer noch hier.",
	"Na, machst du wieder Quatsch?",
	"Wow, du bringst echt frischen Wind rein – im wahrsten Sinne.",
	"Tür auf und Stimmung steigt – fast wie bei ner WG-Party.",
	"Hab dich vermisst. Kurz. Vielleicht.",
	"Das Timing? Chefkochwürdig.",
	"Na, wieder auf geheimer Mission?",
	"Du hast diesen Raum deutlich verbessert. Optisch. Geräuschlich eher weniger.",
	"Wenn Blicke kühlen könnten, wär ich arbeitslos.",
	"Sag nix – ich fühl’s schon.",
	"Heute schon genug Chaos angerichtet, oder kommt noch was?",
	"Du bringst definitiv Energie… oder Hunger. Irgendwas davon.",
	"Nichts sagt „Ich lebe“ so wie du mit dieser Tür.",
	"Ah, der Meister der spontanen Entscheidungen!",
}

func _() {

	apiKey := ""
	client := openai.NewClient(apiKey, "")
	// TTS
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // Maximal 5 parallele TTS-Jobs (kannst du anpassen)

	for i, greeting := range greetings {
		wg.Add(1)
		go func(i int, greeting string) {
			defer wg.Done()

			sem <- struct{}{} // blockiert, wenn mehr als 5 aktiv sind
			defer func() { <-sem }()

			audio, err := client.TextToSpeech(greeting)
			if err != nil {
				fmt.Printf("Error generating audio for greeting '%s': %v\n", greeting, err)
				return
			}

			filename := fmt.Sprintf("greeting_%d.mp3", i+1)
			if err := os.WriteFile(filename, audio, 0644); err != nil {
				fmt.Printf("Error saving audio to file '%s': %v\n", filename, err)
			} else {
				fmt.Printf("Audio saved to %s\n", filename)
			}
		}(i, greeting)
	}

	wg.Wait()
	fmt.Println("Alle verarbeitet.")
}
