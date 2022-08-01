package drill

import (
	mediaConvertAccessor "impruviService/accessor/mediaconvert"
	drillDao "impruviService/dao/drill"
	"log"
)

func CreateDrill(drill *drillDao.DrillDB) (*drillDao.DrillDB, error) {
	drill, err := drillDao.CreateDrill(drill)
	if err != nil {
		return nil, err
	}

	err = startMediaConversion(drill, nil)
	if err != nil {
		return nil, err
	}
	return drill, nil
}

func UpdateDrill(drill *drillDao.DrillDB) error {
	oldDrill, err := drillDao.GetDrillById(drill.DrillId)
	if err != nil {
		log.Printf("Error while getting drill by id: %v. Error: %v\n", drill.DrillId, err)
		return err
	}
	err = drillDao.PutDrill(drill)
	if err != nil {
		log.Printf("Error while putting drill: %+v. Error: %v\n", drill, err)
		return err
	}

	return startMediaConversion(drill, oldDrill)
}

func DeleteDrill(drillId string) error {
	return drillDao.DeleteDrill(drillId)
}

func startMediaConversion(newDrill, oldDrill *drillDao.DrillDB) error {
	if shouldConvertFrontVideo(newDrill, oldDrill) {
		err := mediaConvertAccessor.StartJob(newDrill.Demos.Front.FileLocation, &mediaConvertAccessor.Metadata{
			Type: mediaConvertAccessor.DemoVideo,
			DemoVideoMedata: mediaConvertAccessor.DemoVideoMedata{
				DrillId: newDrill.DrillId,
				Angle:   string(drillDao.FrontAngle),
			},
		})
		if err != nil {
			log.Printf("Error while starting conversion for front angle on drill: %+v. Error: %v\n", newDrill, err)
			return err
		}
	}

	if shouldConvertSideVideo(newDrill, oldDrill) {
		err := mediaConvertAccessor.StartJob(newDrill.Demos.Side.FileLocation, &mediaConvertAccessor.Metadata{
			Type: mediaConvertAccessor.DemoVideo,
			DemoVideoMedata: mediaConvertAccessor.DemoVideoMedata{
				DrillId: newDrill.DrillId,
				Angle:   string(drillDao.SideAngle),
			},
		})
		if err != nil {
			log.Printf("Error while starting conversion for side angle on drill: %+v. Error: %v\n", newDrill, err)
			return err
		}
	}

	if shouldConvertCloseVideo(newDrill, oldDrill) {
		err := mediaConvertAccessor.StartJob(newDrill.Demos.Close.FileLocation, &mediaConvertAccessor.Metadata{
			Type: mediaConvertAccessor.DemoVideo,
			DemoVideoMedata: mediaConvertAccessor.DemoVideoMedata{
				DrillId: newDrill.DrillId,
				Angle:   string(drillDao.CloseAngle),
			},
		})
		if err != nil {
			log.Printf("Error while starting conversion for close angle on drill: %+v. Error: %v\n", newDrill, err)
			return err
		}
	}

	return nil
}

func shouldConvertFrontVideo(newDrill, oldDrill *drillDao.DrillDB) bool {
	if newDrill.Demos == nil || newDrill.Demos.Front == nil {
		return false
	}

	return oldDrill == nil || oldDrill.Demos == nil || oldDrill.Demos.Front == nil || oldDrill.Demos.Front.FileLocation != newDrill.Demos.Front.FileLocation
}

func shouldConvertSideVideo(newDrill, oldDrill *drillDao.DrillDB) bool {
	if newDrill.Demos == nil || newDrill.Demos.Side == nil {
		return false
	}

	return oldDrill == nil || oldDrill.Demos == nil || oldDrill.Demos.Side == nil || oldDrill.Demos.Side.FileLocation != newDrill.Demos.Side.FileLocation
}

func shouldConvertCloseVideo(newDrill, oldDrill *drillDao.DrillDB) bool {
	if newDrill.Demos == nil || newDrill.Demos.Close == nil {
		return false
	}

	return oldDrill == nil || oldDrill.Demos == nil || oldDrill.Demos.Close == nil || oldDrill.Demos.Close.FileLocation != newDrill.Demos.Close.FileLocation
}
