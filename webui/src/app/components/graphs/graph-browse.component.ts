import { Component, OnInit, AfterViewInit, ViewChild } from '@angular/core';
import { GraphBean, NodeBean, GraphVis } from '../../models/graph/graph-bean';
import { Observable } from 'rxjs';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { GraphsStoreService } from './../../stores/graphs-store.service';
import { CytoGraphComponent } from '../../widget/cytograph/cytograph.component';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-graph-browse',
  templateUrl: './graph-browse.component.html',
  styleUrls: [],
})

@AutoUnsubscribe()
export class GraphBrowseComponent implements OnInit {

  @ViewChild('cytograph') cytograph: CytoGraphComponent;
  protected deploymentStream: Observable<GraphBean>;
  private api: any;

  constructor(
    private graphsStoreService: GraphsStoreService,
  ) {
    this.deploymentStream = this.graphsStoreService.deployments();
  }

  private _graphData: any = {
    nodes: [],
    edges: []
  };

  ngOnInit() {
    this.deploymentStream.subscribe(
      (graph: GraphBean) => {
        const dejaVuDomains = {};
        const dejaVuEnvironments = {};
        const dejaVuNodes = {};
        if (graph.nodes === undefined || graph.edges === undefined) {
          return;
        }
        graph.nodes.forEach((node, index) => {
          if (index > 10000) {
            return;
          }
          dejaVuNodes[node.id] = true;
          const env = node.properties.environment.slug;
          const domain = env + '/' + node.properties.application.domain;
          if (dejaVuEnvironments[env] === undefined) {
            this._graphData.nodes.push({
              classes: 'environment',
              data: {
                id: env,
                type: 'environment',
                color: node.properties.environment.properties.color  || 'gray',
                name: env,
              }
            });
            dejaVuEnvironments[env] = true;
          }
          if (dejaVuDomains[domain] === undefined) {
            this._graphData.nodes.push({
              classes: 'domain',
              data: {
                id: domain,
                name: node.properties.application.domain,
                color: node.properties.environment.properties.color || 'gray',
                type: 'domain',
                parent: env,
              }
            });
            dejaVuDomains[domain] = true;
          }
          this._graphData.nodes.push({
            classes: 'application',
            data: {
              id: node.id,
              name: node.name,
              type: 'application',
              domain: node.properties.application.domain,
              parent: domain,
            }
          });
        });
        graph.edges.forEach(edge => {
          if (dejaVuNodes[edge.from] === undefined || dejaVuNodes[edge.to] === undefined) {
            return;
          }
          this._graphData.edges.push({
            data: {
              source: edge.from,
              target: edge.to,
            }
          });
        });
        this.cytograph.load(this._graphData);
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }
}
