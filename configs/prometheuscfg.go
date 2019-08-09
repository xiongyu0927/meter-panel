package configs

import "log"

//GetProConfig4ENV is used for Get prometheus configs from env
func GetProConfig4ENV(K8sconfigs HumanAllK8SConfigs) map[string]string {
	ProCfg := make(map[string]string)
	for k, _ := range K8sconfigs {
		tmp := GetSingleEnvConfigs(k + "Pro")
		if tmp != "" {
			ProCfg[k] = tmp
		} else {
			log.Println("this cluster " + k + " doesn't have the prometheus config, we will try use the podip to access it")
			return nil
		}
	}
	return ProCfg
}
