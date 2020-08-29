import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-toolbar',
  exportAs: 'appToolbar',
  template: '<ng-content></ng-content>',
  styleUrls: ['./toolbar.component.scss'],
})
export class ToolbarComponent implements OnInit {

  constructor() {
  }

  ngOnInit(): void {
  }

}
