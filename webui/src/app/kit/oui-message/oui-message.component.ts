import { Component, OnInit, Output, EventEmitter, Input, ElementRef, OnDestroy, Renderer2 } from '@angular/core';

import { ActivatedRoute } from '@angular/router';
import { ApplicationBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { UiKitMenuItem } from '../../models/kit/navbar';
import { DomHandler } from 'primeng/primeng';

@Component({
  selector: 'app-oui-message',
  templateUrl: './oui-message.component.html',
  styleUrls: [],
  providers: [DomHandler]
})
export class OuiMessageComponent implements OnInit {

  _message: String;
  public visible = true;

  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
  ) {
  }

  ngOnInit() {
  }

  @Input() get message(): String {
    return this._message;
  }

  /**
   * hide
   */
  public hide() {
    this.visible = false;
  }

  /**
   * setter
   */
  set message(val: String) {
    this._message = val;
  }

  /**
   * emit selection
   * @param event
   * @param type
   * @param page
   */
  onSelect(event: any, type: string, page: string) {
    event.data = {
      type: type,
    };
    this.select.emit(event);
  }
}
