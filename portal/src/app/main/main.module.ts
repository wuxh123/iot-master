import {NgModule} from '@angular/core';

import {IconsProviderModule} from './icons-provider.module';
import {NzLayoutModule} from 'ng-zorro-antd/layout';
import {NzMenuModule} from 'ng-zorro-antd/menu';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {CommonModule} from '@angular/common';
import {HttpClientModule} from '@angular/common/http';
import {MainRoutingModule} from './main-routing.module';

import {MainComponent} from './main.component';
import {
  NzButtonModule,
  NzCheckboxModule, NzCollapseModule, NzDividerModule, NzDrawerModule,
  NzFormModule,
  NzIconModule,
  NzInputModule, NzInputNumberModule,
  NzModalModule, NzPopconfirmModule, NzSelectModule, NzStatisticModule, NzSwitchModule,
  NzTableModule, NzTabsModule,
  NzToolTipModule
} from 'ng-zorro-antd';
import {DashComponent} from './dash/dash.component';
import {MomentModule} from 'ngx-moment';
import {UiModule} from '../ui/ui.module';
import {TabsComponent} from './tabs/tabs.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {LinkComponent} from './link/link.component';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {PluginComponent} from './plugin/plugin.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';
import {ProjectComponent} from './project/project.component';
import {ProjectEditComponent} from './project-edit/project-edit.component';
import {ProjectDetailComponent} from './project-detail/project-detail.component';
import {TemplateComponent} from './template/template.component';
import {TemplateEditComponent} from './template-edit/template-edit.component';
import {TemplateDetailComponent} from './template-detail/template-detail.component';
import {ElementComponent} from './element/element.component';
import {ElementEditComponent} from './element-edit/element-edit.component';
import {ElementDetailComponent} from './element-detail/element-detail.component';
import {NzSpaceModule} from 'ng-zorro-antd/space';


@NgModule({
  declarations: [MainComponent, TabsComponent, DashComponent,
    TunnelComponent, TunnelEditComponent,
    LinkComponent, LinkEditComponent, LinkMonitorComponent,
    PluginComponent, PluginEditComponent,
    ProjectComponent, ProjectEditComponent, ProjectDetailComponent,
    TemplateComponent, TemplateEditComponent, TemplateDetailComponent,
    ElementComponent, ElementEditComponent, ElementDetailComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MomentModule,
    UiModule,
    // IconsProviderModule,
    // BrowserModule,
    NzIconModule,
    NzLayoutModule,
    NzMenuModule,
    HttpClientModule,
    MainRoutingModule,
    NzToolTipModule,
    NzTableModule,
    NzModalModule,
    NzFormModule,
    NzButtonModule,
    NzInputModule,
    NzCheckboxModule,
    NzSwitchModule,
    NzPopconfirmModule,
    IconsProviderModule,
    NzDividerModule,
    NzDrawerModule,
    NzSelectModule,
    NzSpaceModule,
    NzInputNumberModule,
    NzStatisticModule,
    NzTabsModule,
    NzCollapseModule,
  ],
  bootstrap: [MainComponent]
})
export class MainModule {
}
