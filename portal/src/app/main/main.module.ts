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
  NzCheckboxModule, NzDividerModule, NzDrawerModule,
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
import {ModelComponent} from './model/model/model.component';
import {ModelEditComponent} from './model/model-edit/model-edit.component';
import {ModelVariableComponent} from './model/variable/model-variable.component';
import {ModelVariableEditComponent} from './model/variable-edit/model-variable-edit.component';
import {ModelBatchComponent} from './model/batch/model-batch.component';
import {ModelBatchEditComponent} from './model/batch-edit/model-batch-edit.component';
import {ModelJobComponent} from './model/job/model-job.component';
import {ModelJobEditComponent} from './model/job-edit/model-job-edit.component';
import {ModelStrategyComponent} from './model/strategy/model-strategy.component';
import {ModelStrategyEditComponent} from './model/strategy-edit/model-strategy-edit.component';
import {ModelAdapterComponent} from './model/adapter/model-adapter.component';
import {ModelAdapterEditComponent} from './model/adapter-edit/model-adapter-edit.component';


@NgModule({
  declarations: [MainComponent, TabsComponent, DashComponent,
    TunnelComponent, TunnelEditComponent,
    LinkComponent, LinkEditComponent, LinkMonitorComponent,
    PluginComponent, PluginEditComponent,
    ModelComponent, ModelEditComponent,
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
  ],
  bootstrap: [MainComponent]
})
export class MainModule {
}
