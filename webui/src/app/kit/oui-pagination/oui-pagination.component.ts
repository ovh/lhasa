import { Component, OnInit, Output, EventEmitter, Input, ElementRef, OnDestroy, Renderer2 } from '@angular/core';

import { ActivatedRoute } from '@angular/router';
import { ApplicationBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { UiKitMenuItem } from '../../models/kit/navbar';
import { DomHandler } from 'primeng/primeng';
import { each } from 'lodash';
import { UiKitPaginate } from '../../models/kit/paginate';

class LogicalPage {
  index: number;
  disabled: string;
  selected: string;
}

@Component({
  selector: 'app-oui-pagination',
  templateUrl: './oui-pagination.component.html',
  styleUrls: ['./oui-pagination.component.css'],
  providers: [DomHandler]
})
export class OuiPaginationComponent implements OnInit {

  _metadata: UiKitPaginate;
  pages: LogicalPage[];

  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
  ) {
  }

  ngOnInit() {
  }

  @Input() get metadata(): UiKitPaginate {
    return this._metadata;
  }

  /**
   * setter
   */
  set metadata(val: UiKitPaginate) {
    this._metadata = val;
    this.update();
  }

  /**
   * update style
   */
  update() {
    this.pages = [];
    for (let index = 0; index < this._metadata.totalPages; index++) {
      let status = 'disabled';
      let selected = 'oui-button oui-button_primary oui-button_small-width oui-pagination-button_selected';
      if (index !== this._metadata.number) {
        status = '';
        selected = 'oui-button oui-button_secondary oui-button_small-width';
      }
      this.pages.push({
        index: index,
        disabled: status,
        selected: selected
      });
    }
  }

  /**
   * emit selection
   * @param event
   * @param type
   * @param page
   */
  onSelect(event: any, type: string, page: number) {
    event.data = {
      page: this.compute(this._metadata, type, page),
      metadata: this._metadata
    };
    this.select.emit(event);
  }

  /**
   * refresh metadata
   */
  RefreshMetadata(metadata: UiKitPaginate, type: string, page: number) {
    const event: any = {};
    this._metadata = metadata;
    event.data = {
      page: this.compute(this._metadata, type, page),
      metadata: this._metadata
    };
    this.update();
    this.select.emit(event);
  }

  /**
   * compute next page
   */
  compute(metadata: UiKitPaginate, type: string, page: number): number {
    let go = 0;
    switch (type) {
      case 'next':
        go = metadata.number + 1;
        break;
      case 'previous':
        go = metadata.number - 1;
        break;
      case 'select':
        go = Number(page);
        break;
    }
    if (go >= metadata.totalPages) {
      go = metadata.totalPages - 1;
    }
    if (go < 0) {
      go = 0;
    }
    return go;
  }
}
