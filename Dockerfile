FROM docker.elastic.co/elasticsearch/elasticsearch-oss:7.9.3

RUN elasticsearch-plugin install analysis-kuromoji && \
    elasticsearch-plugin install analysis-icu

USER elasticsearch
