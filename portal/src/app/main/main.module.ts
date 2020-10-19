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
import {TunnelComponent} from './tunnel/tunnel.component';
import {LinkComponent} from './link/link.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {NzSpaceModule} from 'ng-zorro-antd/space';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {PluginComponent} from './plugin/plugin.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';
import {TabsComponent} from './tabs/tabs.component';
import {ProjectComponent} from './project/project.component';
import {ProjectEditComponent} from './project-edit/project-edit.component';
import {ModelVariableComponent} from './project-variable/model-variable.component';
import {ModelVariableEditComponent} from './project-variable-edit/model-variable-edit.component';
import {ModelBatchComponent} from './project-batch/model-batch.component';
import {ModelBatchEditComponent} from './project-batch-edit/model-batch-edit.component';
import {ModelJobComponent} from './project-job/model-job.component';
import {ModelJobEditComponent} from './project-job-edit/model-job-edit.component';
import {ModelStrategyComponent} from './project-strategy/model-strategy.component';
import {ModelStrategyEditComponent} from './project-strategy-edit/model-strategy-edit.component';
import {ModelAdapterComponent} from './project-adapter/model-adapter.component';
import {ModelAdapterEditComponent} from './project-adapter-edit/model-adapter-edit.component';
import {ProjectDetailComponent} from './project-detail/project-detail.component';


@NgModule({
  declarations: [MainComponent, TabsComponent, DashComponent,
    TunnelComponent, TunnelEditComponent,
    LinkComponent, LinkEditComponent, LinkMonitorComponent,
    PluginComponent, PluginEditComponent,
    ProjectComponent, ProjectEditComponent, ProjectDetailComponent,
    ModelAdapterComponent, ModelAdapterEditComponent,
    ModelVariableComponent, ModelVariableEditComponent,
    ModelBatchComponent, ModelBatchEditComponent,
    ModelJobComponent, ModelJobEditComponent,
    ModelStrategyComponent, ModelStrategyEditComponent,
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
