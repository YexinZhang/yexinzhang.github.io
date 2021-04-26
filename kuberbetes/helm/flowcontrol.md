# if
{{ if PIPELINE }}
{{ else if OTHER PEPELINE }}
{{ else }}
{{ end }}


# with
限定访问域
但是这里有个注意事项，在限定的作用域内，无法使用.访问父作用域的对象。
我们可以使用$从父作用域中访问Release.Name对象
{{ with PIPELINE }}
  <!--  restricted scope -->
{{ end }}

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-configmap
data:
{{- with .Values.favorite }}
  maotai: {{ .drink | default "flySky" | quote }}
  kangshifu: {{ .food | upper | quote }}
  release: {{ $.Release.Name }}
{{- end }}
```

# range
{{- range .Values.pizzaToppings }}
    {{ . | title | quote  }}
{{- end }}
  niupi: |-
{{- range $index, $topping := .Values.pizzaToppings }}
    {{ printf "in_%v"  $index }}: {{ $topping }}
{{- end }}

{{- range $key, $val :=  .Values.favorite }}
  {{ $key  }}: {{ $val | quote }}
{{- end }}

# example
```
favorite:
  drink: coffee
  food: qqmian
pizzaToppings:
  - mogu
  - hotdog
  - qingcai
  - peppers
  - onions
```