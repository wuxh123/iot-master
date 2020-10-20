import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ToolbarComponent} from './toolbar/toolbar.component';
import { FormGridComponent } from './form-grid/form-grid.component';
import {NzGridModule} from 'ng-zorro-antd';


@NgModule({
  declarations: [
    ToolbarComponent,
    FormGridComponent,
  ],
  exports: [
    ToolbarComponent,
    FormGridComponent,
  ],
  imports: [
    CommonModule,
    // NzLayoutModule,
    NzGridModule,
  ]
})
export class UiModule {
}
