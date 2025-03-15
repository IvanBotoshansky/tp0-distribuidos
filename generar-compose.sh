#!/bin/bash
python3 mi-generador.py $1 $2

if [ $? -eq 0 ]; then
    echo "Archivo $1 generado exitosamente con $2 clientes"
else
    echo "Error al generar el archivo"
    exit 1
fi