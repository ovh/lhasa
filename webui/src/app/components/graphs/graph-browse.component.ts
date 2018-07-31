import { GraphsStoreService } from './../../stores/graphs-store.service';
import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { find } from 'lodash';

import { GraphBean, GraphVis, NodeBean } from '../../models/graph/graph-bean';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { GraphComponent } from '../../widget/graph/graph.component';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-graph-browse',
  templateUrl: './graph-browse.component.html',
  styleUrls: [],
})
@AutoUnsubscribe()
export class GraphBrowseComponent implements OnInit, AfterViewInit {

  /**
   * internal streams and store
   */
  public display = false;
  public nodes: NodeBean[] = [];
  protected deploymentStream: Observable<GraphBean>;

  public deployments: GraphBean = {
    nodes: [],
    edges: [],
    options: {}
  };

  public deploymentsVis: GraphVis = {
    nodes: [],
    edges: [],
  };

  public deploymentsVisOptions: {};

  @ViewChild('deploymentsGraph') deploymentsGraph: GraphComponent;

  constructor(
    private graphsStoreService: GraphsStoreService
  ) {
    /**
     * subscribe
     */
    this.deploymentStream = this.graphsStoreService.deployments();
  }

  ngOnInit() {
    this.deploymentStream.subscribe(
      (graph: GraphBean) => {
        this.deployments = graph;
        if (graph.nodes && graph.edges) {
          // Compute data
          this.deploymentsVis.nodes = [];
          graph.nodes.forEach(node => {
            this.deploymentsVis.nodes.push({
              id: node.id,
              label: node.name,
              group: node.properties.environment.slug,
              environment: node.properties.environment.slug,
              domain: node.properties.application.domain,
              application: node.properties.application.name,
            });
          });
          // Compute data
          this.deploymentsVis.edges = [];
          graph.edges.forEach(edge => {
            this.deploymentsVis.edges.push({
              id: edge.id,
              from: edge.from,
              to: edge.to,
              label: edge.type
            });
          });
        }
        // Compute data
        this.deploymentsVisOptions = this.deployments.options;
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }

  ngAfterViewInit() {
  }

  /**
   * event handler
   */
  public onEvent(event: any) {
    if (event.type === 'selectNode') {
      const nodes = <string[]> event.nodes;
      this.nodes = [];
      nodes.forEach(id => {
        this.nodes.push(find(this.deployments.nodes, (item: NodeBean) => {
          return item.id === id;
        }));
      });
    }
    if (event.type === 'doubleClick') {
      this.display = true;
    }
  }

  /**
   * pretty
   */
  public pretty(data: any): string {
    return JSON.stringify(data, null, 2);
  }
}
