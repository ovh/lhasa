import { GraphsStoreService } from './../../stores/graphs-store.service';
import { DataGraphService } from './../../services/data-graph.service';
import { Component, OnInit, AfterViewInit, ViewChild } from '@angular/core';
import { ApplicationsStoreService, LoadApplicationsAction, SelectApplicationAction } from '../../stores/applications-store.service';
import { Store } from '@ngrx/store';
import { ApplicationBean, DeploymentBean, DomainBean } from '../../models/commons/applications-bean';
import { DataApplicationService } from '../../services/data-application-version.service';
import { ContentListResponse } from '../../models/commons/entity-bean';
import { find } from 'lodash';

import { GraphBean, NodeBean, GraphVis } from '../../models/graph/graph-bean';
import { DataDeploymentService } from '../../services/data-deployment.service';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { GraphComponent } from '../../widget/graph/graph.component';

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
  protected deploymentStream: Store<GraphBean>;

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
        // Compute data
        this.deploymentsVis.nodes = [];
        graph.nodes.forEach(node => {
          this.deploymentsVis.nodes.push({
            id: node.id,
            label: node.name,
            group: node.type
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
