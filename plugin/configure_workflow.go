package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
)

func workflowCreateJSON(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg workflowConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return err
	}

	if err := cfg.createWorkflow(); err != nil {
		return err
	}

	return nil
}

func postContainerConfig(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v", err)
		return
	}

	cfg, err := parseWebContainerConfig(r.Form)
	if err != nil {
		log.Println(err)
		return
	}

	JSON, err := json.Marshal(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	fname := cfg.WorkflowName + ".json"
	if err := os.WriteFile(fname, JSON, 0644); err != nil {
		log.Println(err)
		return
	}

	if err := cfg.createWorkflow(); err != nil {
		log.Println(err)
		return
	}

	tmpl := template.Must(template.New("review").Parse(reviewHTML))
	if err := tmpl.Execute(w, cfg); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return
		} else {
			log.Println(err)
			return
		}
	}
}

func workflowCreateWeb() error {
	server := &http.Server{
		Addr:    ":5000",
		Handler: nil,
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index").Parse(indexHTML))
		if err := tmpl.Execute(rw, nil); err != nil {
			if errors.Is(err, syscall.EPIPE) {
				return
			} else {
				log.Println(err)
				server.Close()
			}
		}
	})
	http.HandleFunc("/post", postContainerConfig)
	http.HandleFunc("/quit", func(rw http.ResponseWriter, r *http.Request) { server.Close() })

	fmt.Println("Navigate to 'localhost:5000' to setup your workflow")
	if err := server.ListenAndServe(); err == http.ErrServerClosed {
		fmt.Println("Server closed")
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func parseWebContainerConfig(formData map[string][]string) (workflowConfig, error) {
	var cfg workflowConfig

	pow := func(n, m int64) int64 {
		if m == 0 {
			return 1
		}
		result := n
		for i := int64(2); i <= m; i++ {
			result *= n
		}
		return result
	}

	if val, ok := formData["workflowName"]; ok {
		if len(val) > 0 {
			cfg.WorkflowName = val[0]
		} else {
			return cfg, fmt.Errorf("error parsing workflow name")
		}
	} else {
		return cfg, fmt.Errorf("error parsing workflow name")
	}

	if val, ok := formData["applicationContainer.name"]; ok {
		if len(val) > 0 {
			cfg.ApplicationContainer.Name = val[0]
		} else {
			return cfg, fmt.Errorf("error parsing application container name")

		}
	} else {
		return cfg, fmt.Errorf("error parsing application container name")
	}

	if val, ok := formData["applicationContainer.inPath"]; ok {
		if len(val) > 0 {
			cfg.ApplicationContainer.InPath = val[0]
		} else {
			return cfg, fmt.Errorf("error parsing application input path")
		}
	} else {
		return cfg, fmt.Errorf("error parsing application input path")
	}

	inName, ok1 := formData["inputContainer.name"]
	ininPath, ok2 := formData["inputContainer.inPath"]
	inSize, ok3 := formData["inputContainer.size"]
	inSizeUnit, ok4 := formData["inputContainer.sizeUnit"]

	if ok1 && ok2 && ok3 && ok4 &&
		(len(inName) == len(ininPath) && len(inName) == len(inSize) && len(inName) == len(inSizeUnit)) {
		for i := 0; i < len(inName); i++ {
			size, err := strconv.ParseInt(inSize[i], 0, 64)
			if err != nil {
				return cfg, fmt.Errorf("error parsing input containers")
			}
			sizeUnit, err := strconv.ParseInt(inSizeUnit[i], 0, 64)
			if err != nil {
				return cfg, fmt.Errorf("error parsing input containers")
			}
			size *= pow(1024, sizeUnit)
			cfg.InputContainer = append(cfg.InputContainer, containerConfig{
				Name:   inName[i],
				InPath: ininPath[i],
				Size:   size,
			})
		}
	} else {
		return cfg, fmt.Errorf("error parsing input containers")
	}

	if val, ok := formData["outputContainer.name"]; ok {
		if len(val) > 0 {
			cfg.OutputContainer.Name = val[0]
		} else {
			return cfg, fmt.Errorf("error parsing output container name")
		}
	} else {
		return cfg, fmt.Errorf("error parsing output container name")
	}

	outSize, ok1 := formData["outputContainer.size"]
	outSizeUnit, ok2 := formData["outputContainer.sizeUnit"]

	if ok1 && ok2 &&
		(len(outSize) > 0) && (len(outSizeUnit) > 0) {
		size, err := strconv.ParseInt(outSize[0], 0, 64)
		if err != nil {
			return cfg, fmt.Errorf("error parsing output container")
		}
		sizeUnit, err := strconv.ParseInt(outSizeUnit[0], 0, 64)
		if err != nil {
			return cfg, fmt.Errorf("error parsing output container")
		}
		size *= pow(1024, sizeUnit)
		cfg.OutputContainer.Size = size
	} else {
		return cfg, fmt.Errorf("error parsing output container")
	}

	return cfg, nil
}
