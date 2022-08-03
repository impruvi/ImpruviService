package emailtemplate

import (
	"fmt"
	coachDao "impruviService/dao/coach"
	playerFacade "impruviService/facade/player"
	"log"
	"strings"
)

func GetSubscriptionCreatedCoachEmailHtml(player *playerFacade.Player, coach *coachDao.CoachDB) string {
	return fmt.Sprintf(`
<div>%s,</div>
<br />

<div><b>%s %s</b> has just subscribed to train with you!  You have <b>24 hours</b> to build them a plan on the Impruvi app.</div>
<br />

<div>Use their questionnaire results below to help build their plan.</div>
<br />

<div>You can further customize %s’s plan based on their submission videos for your intro session.</div>
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
		player.FirstName,
		player.LastName,
		player.Position,
		player.AgeRange,
		getEquipmentListDisplayValue(player.AvailableEquipment),
		getTrainingLocationListDisplayValue(player.AvailableTrainingLocations),
		player.ShortTermGoal,
		player.LongTermGoal)
}

func GetSubscriptionCreatedCoachEmailText(player *playerFacade.Player, coach *coachDao.CoachDB) string {
	return fmt.Sprintf(`
%s,
 
%s %s has just subscribed to train with you!  You have 24 hours to build them a plan on the Impruvi app.
 
Use their questionnaire results below to help build their plan.
 
You can further customize %s’s plan based on their submission videos for your intro session.
 
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
		player.FirstName,
		player.LastName,
		player.Position,
		player.AgeRange,
		getEquipmentListDisplayValue(player.AvailableEquipment),
		getTrainingLocationListDisplayValue(player.AvailableTrainingLocations),
		player.ShortTermGoal,
		player.LongTermGoal)
}

func GetSubscriptionCreatedPlayerEmailHtml(player *playerFacade.Player, coach *coachDao.CoachDB) string {
	return fmt.Sprintf(`
<div>%s %s,</div>
<br />
 
<div>
Congratulations on starting your first subscription plan with imprüvi! 
Coach %s %s is looking forward to working with you to develop your skills.
</div>
<br />
 
<div>
To start training, download the imprüvi app on the <a href="https://apps.apple.com/us/app/impruvi/id1627911060">App Store</a>. 
To login on the app, use the same email and password you used to create your account on the website.
</div>
<br />

<div>
We are so grateful and excited that you are part of our company’s early development stages. 
Our founding vision was to create a platform that helps players like you achieve their goals. 
We strive to create an impactful experience every step of the way, so please don’t hesitate to 
reach out if you have any questions, concerns or ideas that would make your experience better. 
</div>
<br />

<div>
Reach us anytime at ryan@impruviapp.com or 720-233-1012. 
</div>
<br />
 
<div>
Imprüvi Founders, 
</div>
<br />
 
<div>
Ryan Crowley and John Magnus
</div>
`,
		player.FirstName,
		player.LastName,
		coach.FirstName,
		coach.LastName)
}

func GetSubscriptionCreatedPlayerEmailText(player *playerFacade.Player, coach *coachDao.CoachDB) string {
	return fmt.Sprintf(`
%s %s, 
 
Congratulations on starting your first subscription plan with imprüvi! Coach %s %s is looking forward to working with you to develop your skills. 
 
To start training, download the imprüvi app on the App Store. To login on the app, use the same email and password you used to create your account on the website.
 
We are so grateful and excited that you are part of our company’s early development stages. Our founding vision was to create a platform that helps players like you achieve their goals. We strive to create an impactful experience every step of the way, so please don’t hesitate to reach out if you have any questions, concerns or ideas that would make your experience better. 
 
Reach us anytime at ryan@impruviapp.com or 720-233-1012. 

 
Imprüvi Founders, 
 
Ryan Crowley and John Magnus
`,
		player.FirstName,
		player.LastName,
		coach.FirstName,
		coach.LastName)
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
