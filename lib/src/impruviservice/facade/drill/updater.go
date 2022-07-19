package drill

import (
	drillDao "impruviService/dao/drill"
)

func CreateDrill(drill *drillDao.DrillDB) (*drillDao.DrillDB, error) {
	return drillDao.CreateDrill(drill)
}

func UpdateDrill(drill *drillDao.DrillDB) error {
	return drillDao.PutDrill(drill)
}

func DeleteDrill(drillId string) error {
	return drillDao.DeleteDrill(drillId)
}
