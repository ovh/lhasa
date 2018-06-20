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
  _stack: any;
  _type: String;
  _klass: String;
  public visible = true;

  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
  ) {
    this._type = 'info';
    this.validate();
  }

  ngOnInit() {
  }

  @Input() get message(): String {
    return this._message;
  }

  @Input() get type(): String {
    return this._type;
  }

  @Input() get stack(): any {
    return this._stack;
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
  set type(val: String) {
    this._type = val;
    this.validate();
  }
  set stack(val: any) {
    this._stack = val;
  }

  validate() {
    switch(this._type) {
      case 'info':
        this._klass = 'oui-message oui-message_info';
      break;
      case 'error':
        this._klass = 'oui-message oui-message_error';
      break;
    }
  }

  /**
   * emit selection
   * @param event
   * @param type
   * @param page
   */
  onSelect(event: any, type: string, page: string) {
    event.data = {
      sender: this,
      message: this.message,
      type: type,
    };
    this.select.emit(event);
  }
}
