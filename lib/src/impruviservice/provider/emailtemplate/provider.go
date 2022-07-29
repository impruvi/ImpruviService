package emailtemplate

import (
	"fmt"
	coachDao "impruviService/dao/coach"
	playerFacade "impruviService/facade/player"
	"log"
	"strings"
)

func GetSubscriptionCreateHtml(player *playerFacade.Player, coach *coachDao.CoachDB) string {
	return fmt.Sprintf(`
<div>%s,</div>
<br />

<div><b>%s %s</b> has just subscribed to train with you!  You have <b>24 hours</b> to build them a plan on the Impruvi app.</div>
<br />

<div>Use their questionnaire results below to help build their plan.</div>
<br />

<div>You can further customize Christian’s plan based on their submission videos for your intro session.</div>
<br />

<div>How to manage the coach’s portal:</div>
<div>link to coaches onboarding doc</div>
<br />

<div><b>Name: </b>%s %s</div>
<div><b>Position: </b>%s</div>
<div><b>Age </b>%s</div>
<div><b>Equipment: </b>%s</div>
<div><b>Where they will train: </b>%s</div>
<div><b>Monthly goals: </b>%s</div>
<div><b>Long term goals: </b>%s</div>
`,
		coach.FirstName,
		player.FirstName,
		player.LastName,
		player.FirstName,
		player.LastName,
		player.Position,
		player.AgeRange,
		getEquipmentListDisplayValue(player.AvailableEquipment),
		getTrainingLocationListDisplayValue(player.AvailableTrainingLocations),
		player.ShortTermGoal,
		player.LongTermGoal)
}

func GetSubscriptionCreatedText(player *playerFacade.Player, coach *coachDao.CoachDB) string {
	return fmt.Sprintf(`
%s,
 
%s %s has just subscribed to train with you!  You have 24 hours to build them a plan on the Impruvi app.
 
Use their questionnaire results below to help build their plan.
 
You can further customize Christian’s plan based on their submission videos for your intro session.
 
How to manage the coach’s portal:
<link to coaches onboarding doc>
 
Name: %s %s
Position: %s
Age %s
Equipment: %s
Where they will train: %s
Monthly goals: %s
Long term goals: %s
`,
		coach.FirstName,
		player.FirstName,
		player.LastName,
		player.FirstName,
		player.LastName,
		player.Position,
		player.AgeRange,
		getEquipmentListDisplayValue(player.AvailableEquipment),
		getTrainingLocationListDisplayValue(player.AvailableTrainingLocations),
		player.ShortTermGoal,
		player.LongTermGoal)
}

func getTrainingLocationListDisplayValue(trainingLocations []string) string {
	return strings.Join(trainingLocations, ", ")
}

func getEquipmentListDisplayValue(equipmentTypeList []string) string {
	displayValues := make([]string, 0)
	for _, equipmentType := range equipmentTypeList {
		displayValue := getEquipmentDisplayValue(equipmentType)
		if displayValue != "" {
			displayValues = append(displayValues, displayValue)
		} else {
			log.Printf("[ERROR] unexpected equipment type: %v\n", equipmentType)
		}
	}

	return strings.Join(displayValues, ", ")
}

func getEquipmentDisplayValue(equipmentType string) string {
	if equipmentType == "BALL" {
		return "ball"
	} else if equipmentType == "CONE" {
		return "cone"
	} else if equipmentType == "GOAL" {
		return "goal"
	} else if equipmentType == "AGILITY_LADDER" {
		return "agility ladder"
	}
	return ""
}
